package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 1. 定义服务结构体（和标准 RPC 完全一致）
type Calculator int

// 2. 定义远程方法（严格遵守 RPC 方法签名规则）
// Add：计算两数之和
func (c *Calculator) Add(args [2]int, reply *int) error {
	*reply = args[0] + args[1]
	return nil
}

// Divide：计算两数之商（处理除数为0的异常）
func (c *Calculator) Divide(args [2]int, reply *float64) error {
	if args[1] == 0 {
		return errors.New("除数不能为0")
	}
	*reply = float64(args[0]) / float64(args[1])
	return nil
}

func main_old() {
	// 3. 注册 RPC 服务（和标准 RPC 一致）
	calc := new(Calculator)
	err := rpc.Register(calc)
	if err != nil {
		panic("注册服务失败：" + err.Error())
	}

	// 4. 开启 TCP 监听（端口 1234）
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic("监听端口失败：" + err.Error())
	}
	defer listener.Close()
	println("jsonrpc 服务端已启动，监听端口 1234...")

	// 5. 处理客户端连接（核心区别：用 jsonrpc.ServeConn 替代 rpc.ServeConn）
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("接收连接失败：" + err.Error())
			continue
		}
		println("收到请求：" + conn.RemoteAddr().String())

		// 临时添加：打印客户端发送的原始数据
		go func() {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("读取客户端数据失败：%s", err)
				return
			}
			log.Printf("客户端原始请求：%s", string(buf[:n]))
			// 把读取的数据写回 conn（否则 ServeConn 会读不到）
			_, _ = conn.Write(buf[:n])
		}()

		// 异步处理连接，使用 jsonrpc 协议
		go jsonrpc.ServeConn(conn)
	}
}
