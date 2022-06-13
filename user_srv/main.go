package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"shop-srvs/user_srv/global"
	"shop-srvs/user_srv/utils"
	"syscall"

	"github.com/hashicorp/consul/api"
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
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	registration := new(api.AgentServiceRegistration)
	registration.Name = "user-srv"
	registration.ID = serviceId
	registration.Port = *port
	registration.Tags = []string{"user", "srv"}
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
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
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
