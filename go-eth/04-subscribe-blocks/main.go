package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 本地 Anvil WebSocket 地址（固定）
const wsURL = "ws://127.0.0.1:8545"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 连接 WebSocket
	fmt.Println("Connecting to Anvil WebSocket:", wsURL)
	client, err := ethclient.DialContext(ctx, wsURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	// 订阅新区块
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalf("订阅失败: %v", err)
	}

	fmt.Println("✅ 订阅成功！正在监听新区块...")
	fmt.Println("按 Ctrl+C 退出\n")

	// 监听退出信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 循环处理事件
	for {
		select {
		case h := <-headers:
			fmt.Printf(
				"🟢 新块 | 高度: %-8d | 哈希: %s\n",
				h.Number.Uint64(),
				h.Hash().Hex(),
			)

		case err := <-sub.Err():
			log.Fatalf("订阅错误: %v", err)

		case sig := <-sigCh:
			fmt.Printf("\n📌 收到信号 %s，程序退出\n", sig)
			return

		case <-ctx.Done():
			fmt.Println("程序退出")
			return
		}
	}
}
