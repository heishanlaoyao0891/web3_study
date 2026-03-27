package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const localRPC = "http://127.0.0.1:8545"
const senderPrivkey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func main() {
	txHashHex := flag.String("tx", "", "transaction hash (for query mode)")
	sendMode := flag.Bool("send", false, "enable send transaction mode")
	toAddrHex := flag.String("to", "", "recipient address (required for send mode)")
	amountEth := flag.Float64("amount", 0, "amount in ETH (required for send mode)")
	flag.Parse()

	//判断操作模式
	if *sendMode {
		if *toAddrHex == "" || *amountEth == 0 {
			log.Fatal("send mode requires --to and --amount flags")
		}
		sendTransaction(*toAddrHex, *amountEth)
	} else {
		if *txHashHex == "" {
			log.Fatal("query mode requires --tx flag, or use --send for send mode")
		}
		queryTransaction(*txHashHex)
	}

}

// 查询交易
func queryTransaction(hash string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, localRPC)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	txHash := common.HexToHash(hash)
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}
	fmt.Println("=== Transaction ===")

	printTxBasicInfo(tx, isPending)

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Printf("receipt unavailable (pending): %v", err)
		return
	}

	fmt.Println("=== Receipt ===")
	printReceiptInfo(receipt)

}
func printTxBasicInfo(tx *types.Transaction, isPending bool) {
	fmt.Printf("Hash        : %s\n", tx.Hash().Hex())
	fmt.Printf("Nonce       : %d\n", tx.Nonce())
	fmt.Printf("Gas         : %d\n", tx.Gas())
	fmt.Printf("To          : %v\n", tx.To())
	fmt.Printf("Value       : %s ETH\n", weiToEth(tx.Value()))
	fmt.Printf("Pending     : %v\n", isPending)
}

func printReceiptInfo(r *types.Receipt) {
	fmt.Printf("状态        : %d (1=成功)\n", r.Status)
	fmt.Printf("区块高度    : %d\n", r.BlockNumber.Uint64())
	fmt.Printf("Gas Used    : %d\n", r.GasUsed)
	fmt.Printf("日志数量    : %d\n", len(r.Logs))
}

// 发送交易
func sendTransaction(to string, amount float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, localRPC)
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer client.Close()
	//解析内置私钥 anvil 第一个账号
	privKey, err := crypto.HexToECDSA(senderPrivkey)
	if err != nil {
		log.Fatalf("invalid privkey: %v", err)
	}
	//发送方地址
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	fromAddr := crypto.PubkeyToAddress(*pubKey)
	toAddr := common.HexToAddress(to)

	//链
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("chain id: %v", err)
	}
	//nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		log.Fatalf("nonce: %v", err)
	}
	//gas 相关
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Fatalf("gas tip: %v", err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("header: %v", err)
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee = big.NewInt(0)
	}

	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	gasLimit := uint64(21000)

	// 金额转换 ETH → Wei
	amountWei := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18))
	valueWei, _ := amountWei.Int(nil)

	// 余额检查
	balance, _ := client.BalanceAt(ctx, fromAddr, nil)
	fmt.Printf("> 发送方余额: %s ETH\n", weiToEth(balance))

	// 构造交易
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     valueWei,
		Data:      nil,
	})

	// 签名
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privKey)
	if err != nil {
		log.Fatalf("sign failed: %v", err)
	}

	// 发送
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatalf("send failed: %v", err)
	}

	// 输出结果
	fmt.Println("\n=== 交易已发送 ===")
	fmt.Printf("From    : %s\n", fromAddr.Hex())
	fmt.Printf("To      : %s\n", toAddr.Hex())
	fmt.Printf("Amount  : %.6f ETH\n", amount)
	fmt.Printf("Tx Hash : %s\n", signedTx.Hash().Hex())
	fmt.Println("\n查询命令：")
	fmt.Printf("go run main.go --tx %s\n", signedTx.Hash().Hex())
}

// wei → ETH 显示用
func weiToEth(wei *big.Int) string {
	f := new(big.Float).SetInt(wei)
	eth := new(big.Float).Quo(f, big.NewFloat(1e18))
	return fmt.Sprintf("%.6f", eth)
}

// 去掉 0x 前缀
func trim0x(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

/**
go run main.go --send --to 0x70997970C51810dc519abffe0E13d384DD20f4b2 --amount 0.5
> 发送方余额: 9999.999113 ETH

=== 交易已发送 ===
From    : 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
To      : 0x70997970c51810DC519aBffe0E13d384dD20f4B2
Amount  : 0.500000 ETH
Tx Hash : 0x9dd1974ad73c30d1259954f15e5d2a6968ad90210d02544091e6c2685f70f098


go run main.go --tx 0x9dd1974ad73c30d1259954f15e5d2a6968ad90210d02544091e6c2685f70f098
=== Transaction ===
Hash        : 0x9dd1974ad73c30d1259954f15e5d2a6968ad90210d02544091e6c2685f70f098
Nonce       : 2
Gas         : 21000
To          : 0x70997970c51810DC519aBffe0E13d384dD20f4B2
Value       : 0.500000 ETH
Pending     : false
=== Receipt ===
状态        : 1 (1=成功)
区块高度    : 3
Gas Used    : 21000
日志数量    : 0
*/
