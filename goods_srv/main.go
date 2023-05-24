package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"srvs/goods_srv/global"
	"srvs/goods_srv/utils"
	"srvs/goods_srv/utils/register/consul"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"srvs/goods_srv/handler"
	"srvs/goods_srv/initialize"
	"srvs/goods_srv/proto/proto"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "IP地址")
	port := flag.Int("port", 0, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

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

	serviceId := fmt.Sprintf("%s", uuid.NewV4())

	registryClient := consul.NewRegistryClient(serverInfo.ConsulInfo.Host, serverInfo.ConsulInfo.Port)

	err = registryClient.Register(serverInfo.Host, *port, serverInfo.Name, serverInfo.Tags, serviceId)

	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
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
	if err = registryClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
