package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"shop-srvs/user_srv/handler"
	"shop-srvs/user_srv/proto/proto"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "IP地址")
	port := flag.Int("port", 50051, "端口号")

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", *ip, *port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(lis)

	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
