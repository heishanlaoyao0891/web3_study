package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 你的Sepolia节点地址
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/eHJVP7yjzEB_1fq6yExoV")
	if err != nil {
		log.Fatal("节点连接失败：", err) // 失败则打印错误原因
	}
	defer client.Close()

	// 获取最新区块号，验证节点可正常查询数据
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal("获取区块信息失败：", err)
	}
	fmt.Printf("节点可用，Sepolia最新区块高度：%d\n", header.Number.Uint64())
}
