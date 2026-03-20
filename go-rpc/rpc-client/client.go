package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

func main_old() {
	// 1. 连接 jsonrpc 服务端（核心区别：用 jsonrpc.Dial 替代 rpc.Dial）
	client, err := jsonrpc.Dial("tcp", "124.223.6.26:9999")
	if err != nil {
		panic("连接服务端失败：" + err.Error())
	}
	defer client.Close()

	// 2. 调用远程方法 Add
	var addResult int
	err = client.Call("Calculator.Add", [2]int{10, 20}, &addResult)
	if err != nil {
		panic("调用 Add 失败：" + err.Error())
	}
	fmt.Printf("10 + 20 = %d\n", addResult) // 输出：10 + 20 = 30

	// 3. 调用远程方法 Divide（正常情况）
	var divResult float64
	err = client.Call("Calculator.Divide", [2]int{20, 5}, &divResult)
	if err != nil {
		panic("调用 Divide 失败：" + err.Error())
	}
	fmt.Printf("20 / 5 = %.2f\n", divResult) // 输出：20 / 5 = 4.00

	// 4. 调用 Divide（异常情况：除数为0）
	err = client.Call("Calculator.Divide", [2]int{10, 0}, &divResult)
	if err != nil {
		fmt.Printf("调用 Divide 失败：%s\n", err.Error()) // 输出：调用 Divide 失败：除数不能为0
	}
}
