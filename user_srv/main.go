package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"shop-srvs/user_srv/handler"
	"shop-srvs/user_srv/initialize"
	"shop-srvs/user_srv/proto/proto"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "IP地址")
	port := flag.Int("port", 50051, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	zap.S().Info("服务启动中...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%v", *ip, *port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(lis)

	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
