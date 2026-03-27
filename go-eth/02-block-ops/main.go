package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	_ "os"
	"time"

	_ "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 本地 Anvil 链 RPC 地址
const localRPC = "http://127.0.0.1:8545"

func main() {
	// 1. 定义命令行参数
	blockNumberFlag := flag.Uint64("number", 0, "查询指定区块号 (0 表示查询最新区块)")
	rangeStartFlag := flag.Uint64("range-start", 0, "批量查询起始区块")
	rangeEndFlag := flag.Uint64("range-end", 0, "批量查询结束区块")
	rateLimitFlag := flag.Int("rate-limit", 200, "批量查询间隔（毫秒）")

	// 解析命令行参数
	flag.Parse()

	// 2. 连接本地 Anvil 链
	fmt.Println("🔗 正在连接本地 Anvil 链：", localRPC)
	client, err := ethclient.Dial(localRPC)
	if err != nil {
		log.Fatalf("❌ 连接 RPC 失败: %v", err)
	}
	defer client.Close()
	fmt.Println("✅ 成功连接本地区块链\n")

	// 3. 执行不同查询逻辑
	if *rangeStartFlag > 0 && *rangeEndFlag > 0 {
		// 批量查询区块范围
		fmt.Printf("📦 开始批量查询区块范围: [%d, %d]\n", *rangeStartFlag, *rangeEndFlag)
		queryBlockRange(client, *rangeStartFlag, *rangeEndFlag, *rateLimitFlag)
	} else if *blockNumberFlag > 0 {
		// 查询单个指定区块
		fmt.Printf("📦 正在查询区块: %d\n", *blockNumberFlag)
		querySingleBlock(client, new(big.Int).SetUint64(*blockNumberFlag))
	} else {
		// 查询最新区块
		fmt.Println("📦 正在查询最新区块...")
		queryLatestBlock(client)
	}
}

// 查询最新区块
func queryLatestBlock(client *ethclient.Client) {
	// 获取最新区块头
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("获取最新区块头失败: %v", err)
	}

	// 获取完整区块
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("获取最新区块失败: %v", err)
	}

	printBlockInfo(block, header.Number.Uint64())
}

// 查询单个指定区块
func querySingleBlock(client *ethclient.Client, number *big.Int) {
	block, err := client.BlockByNumber(context.Background(), number)
	if err != nil {
		log.Fatalf("获取区块 #%d 失败: %v", number.Uint64(), err)
	}
	printBlockInfo(block, number.Uint64())
}

// 批量查询区块范围
func queryBlockRange(client *ethclient.Client, start, end uint64, ms int) {
	if start > end {
		log.Fatal("起始区块不能大于结束区块")
	}

	for i := start; i <= end; i++ {
		fmt.Println("----------------------------------------")
		querySingleBlock(client, new(big.Int).SetUint64(i))
		// 请求间隔
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("✅ 批量查询完成！共查询 %d 个区块\n", end-start+1)
}

// 打印区块详细信息
func printBlockInfo(block *types.Block, number uint64) {
	fmt.Printf("📊 区块信息 #%d\n", number)
	fmt.Printf("哈希: %s\n", block.Hash().Hex())
	fmt.Printf("父哈希: %s\n", block.ParentHash().Hex())
	fmt.Printf("区块高度: %d\n", block.Number().Uint64())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("难度: %d\n", block.Difficulty())
	fmt.Printf("Gas上限: %d\n", block.GasLimit())
	fmt.Printf("已用Gas: %d\n", block.GasUsed())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))

	// 打印区块内的交易哈希
	if len(block.Transactions()) > 0 {
		fmt.Println("💱 交易列表:")
		for idx, tx := range block.Transactions() {
			fmt.Printf("  %d: %s\n", idx+1, tx.Hash().Hex())
		}
	} else {
		fmt.Println("💱 交易列表: 无交易")
	}
}
