package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"shop-srvs/goods_srv/utils/register/consul"
	"shop-srvs/user_srv/global"
	"shop-srvs/user_srv/utils"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop-srvs/user_srv/handler"
	"shop-srvs/user_srv/initialize"
	"shop-srvs/user_srv/proto/proto"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "IP地址")
	port := flag.Int("port", 0, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitMysql()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

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

	serverCfg := global.ServerConfig

	registryClient := consul.NewRegistryClient(serverCfg.ConsulInfo.Host, serverCfg.ConsulInfo.Port)

	serviceId := fmt.Sprintf("%s", uuid.NewV4())

	err = registryClient.Register(serverCfg.Host, *port, serverCfg.Name, serverCfg.Tags, serviceId)
	if err != nil {
		zap.S().Panicf("服务注册失败：", err.Error())
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
