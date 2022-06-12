package main

import (
	"flag"
	"fmt"
	"net"
	"shop-srvs/goods_srv/global"
	"shop-srvs/goods_srv/utils"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop-srvs/goods_srv/handler"
	"shop-srvs/goods_srv/initialize"
	"shop-srvs/goods_srv/proto/proto"
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

	consulInfo := global.ServerConfig.ConsulInfo

	//健康检查
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)

	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = fmt.Sprintf("%s", uuid.NewV4())
	registration.Port = *port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)

	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
