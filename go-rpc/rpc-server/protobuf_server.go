package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"rpc-server/pb"

	"google.golang.org/protobuf/proto"
)

// 定义 RPC 服务结构体
type UserService struct{}

// 实现 Login 方法
func (s *UserService) Login(request []byte, reply *[]byte) error {
	// 反序列化请求
	var user pb.User
	if err := proto.Unmarshal(request, &user); err != nil {
		return err
	}

	log.Printf("收到登录请求：%v", user)

	// 处理业务逻辑
	response := &pb.User{
		Name:  "欢迎 " + user.Name,
		Age:   user.Age,
		IsVip: user.IsVip,
		Tags:  user.Tags,
	}

	// 序列化响应
	data, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	*reply = data
	return nil
}

func main() {
	// 注册 RPC 服务
	err := rpc.Register(new(UserService))
	if err != nil {
		panic("注册服务失败：" + err.Error())
	}

	// 开启 TCP 监听
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic("监听端口失败：" + err.Error())
	}
	defer listener.Close()
	log.Println("JSON-RPC 服务端已启动，监听端口 9999...")

	// 处理客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接收连接失败：" + err.Error())
			continue
		}
		log.Println("收到请求：" + conn.RemoteAddr().String())

		// 异步处理连接
		go jsonrpc.ServeConn(conn)
	}
}
