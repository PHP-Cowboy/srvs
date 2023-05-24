package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"srvs/order_srv/global"
	"srvs/order_srv/utils"
	"srvs/order_srv/utils/register/consul"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"srvs/order_srv/handler"
	"srvs/order_srv/initialize"
	"srvs/order_srv/proto/proto"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "IP地址")
	port := flag.Int("port", 50051, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()

	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	if *port == 0 {
		*port = utils.GetFreePort()
	}

	zap.S().Info("服务启动中...")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//服务注册
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	serverInfo := global.ServerConfig

	//生成注册对象
	serviceId := fmt.Sprintf("%s", uuid.NewV4())

	register_client := consul.NewRegistry(serverInfo.ConsulInfo.Host, serverInfo.ConsulInfo.Port)

	err = register_client.Register(serverInfo.Host, *port, serverInfo.Name, serverInfo.Tags, serviceId)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)

		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
