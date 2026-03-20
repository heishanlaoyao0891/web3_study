package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
	"rpc-client/rpc-client/pb"

	"google.golang.org/protobuf/proto"
)

func main() {
	//创建user对象
	user := &pb.User{
		Name:  "sl",
		Age:   10,
		IsVip: true,
		Tags:  []string{"a", "b"},
	}
	fmt.Printf("原始user: %v\n", user)
	//序列化
	//把protobuf对象转换成二进制字节流，用于网络传输
	data, err := proto.Marshal(user)
	if err != nil {
		log.Fatal("序列化失败", err)
	}
	fmt.Printf("序列化后use=%#x\n", data)
	fmt.Println("调用远程rpc接口")
	client, err := jsonrpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("连接服务端失败")
	}
	defer client.Close()
	var response []byte
	err = client.Call("UserService.Login", data, &response)
	if err != nil {
		log.Fatal("调用失败:", err)
	}

	// 反序列化响应
	var result pb.User
	if err := proto.Unmarshal(response, &result); err != nil {
		log.Fatal("反序列化失败:", err)
	}

	fmt.Printf("响应结果: %v\n", result)

}
