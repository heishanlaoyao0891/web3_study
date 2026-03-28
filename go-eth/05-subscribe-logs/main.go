package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 🔥 你自己的 MyERC20 ABI（我已帮你放好）
const erc20ABIJSON = `[
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "name",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "symbol",
        "type": "string"
      },
      {
        "internalType": "uint256",
        "name": "initialSupply",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "recipient",
        "type": "address"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "spender",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "Approval",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "from",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "to",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "Transfer",
    "type": "event"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "internalType": "address",
        "name": "spender",
        "type": "address"
      }
    ],
    "name": "allowance",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "spender",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "approve",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "account",
        "type": "address"
      }
    ],
    "name": "balanceOf",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "decimals",
    "outputs": [
      {
        "internalType": "uint8",
        "name": "",
        "type": "uint8"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "name",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "symbol",
    "outputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "totalSupply",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "to",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "transfer",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "from",
        "type": "address"
      },
      {
        "internalType": "address",
        "name": "to",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "value",
        "type": "uint256"
      }
    ],
    "name": "transferFrom",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "to",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "mint",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]`

var (
	rpcURL = "ws://127.0.0.1:8545" // 🔥 内置地址，不用配置环境变量
)

func main() {
	contractAddr := flag.String("contract", "", "ERC20 合约地址")
	flag.Parse()

	if *contractAddr == "" {
		log.Fatal("必须传入 --contract 合约地址")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n退出程序...")
		cancel()
	}()

	// 启动带重连的日志监听
	runLogReconnect(ctx, *contractAddr)
}

// 带自动重连的日志订阅
func runLogReconnect(ctx context.Context, contractAddr string) {
	parsedABI, _ := abi.JSON(strings.NewReader(erc20ABIJSON))
	contract := common.HexToAddress(contractAddr)
	attempt := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		attempt++
		log.Printf("🔗 第 %d 次连接: %s", attempt, rpcURL)

		client, err := ethclient.DialContext(ctx, rpcURL)
		if err != nil {
			log.Printf("连接失败: %v", err)
			backoff(ctx, attempt)
			continue
		}

		logsCh := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
			Addresses: []common.Address{contract},
		}, logsCh)
		if err != nil {
			log.Printf("订阅失败: %v", err)
			client.Close()
			backoff(ctx, attempt)
			continue
		}

		log.Println("✅ 监听 ERC20 事件成功 (Transfer / Approval)")

		// 监听循环
		for {
			select {
			case logData := <-logsCh:
				parseLogEvent(&logData, parsedABI)
			case err := <-sub.Err():
				log.Printf("连接断开: %v", err)
				client.Close()
				backoff(ctx, attempt)
				goto RECONNECT
			case <-ctx.Done():
				client.Close()
				return
			}
		}

	RECONNECT:
	}
}

// 解析事件（你原来的逻辑完全保留）
func parseLogEvent(vLog *types.Log, parsedABI abi.ABI) {
	if len(vLog.Topics) == 0 {
		return
	}

	eventTopic := vLog.Topics[0]
	var eventName string
	var eventSig abi.Event

	for name, event := range parsedABI.Events {
		if crypto.Keccak256Hash([]byte(event.Sig)) == eventTopic {
			eventName = name
			eventSig = event
			break
		}
	}

	if eventName == "" {
		fmt.Printf("未知事件 | 区块: %d | Tx: %s\n", vLog.BlockNumber, vLog.TxHash.Hex())
		return
	}

	fmt.Printf("=======================================================\n")
	fmt.Printf("✅ 事件: %s | 区块: %d | Tx: %s\n", eventName, vLog.BlockNumber, vLog.TxHash.Hex())

	indexedParamIndex := 0
	for _, input := range eventSig.Inputs {
		if !input.Indexed {
			continue
		}
		topicIdx := 1 + indexedParamIndex
		indexedParamIndex++
		if topicIdx >= len(vLog.Topics) {
			continue
		}

		topic := vLog.Topics[topicIdx]
		fmt.Printf("  📌 %s: ", input.Name)

		switch input.Type.T {
		case abi.AddressTy:
			fmt.Printf("%s\n", common.BytesToAddress(topic.Bytes()).Hex())
		case abi.UintTy, abi.IntTy:
			fmt.Printf("%s\n", new(big.Int).SetBytes(topic.Bytes()).String())
		default:
			fmt.Printf("%s\n", topic.Hex())
		}
	}

	if len(vLog.Data) > 0 {
		values, _ := parsedABI.Unpack(eventName, vLog.Data)
		nonIdx := 0
		for _, input := range eventSig.Inputs {
			if !input.Indexed {
				if nonIdx < len(values) {
					fmt.Printf("  💰 %s: %v\n", input.Name, values[nonIdx])
					nonIdx++
				}
			}
		}
	}

	fmt.Printf("=======================================================\n\n")
}

// 指数退避重连
func backoff(ctx context.Context, attempt int) {
	sec := int(math.Min(60, math.Pow(2, float64(attempt))))
	d := time.Duration(sec) * time.Second
	log.Printf("⏳ %s 后重连...\n", d)

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
