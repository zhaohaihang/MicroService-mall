package mode

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"github.com/zhaohaihang/goods_service/global"
	"github.com/zhaohaihang/goods_service/utils"
	"go.uber.org/zap"

	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// ReleaseMode release模式注册服务

func ReleaseMode(server *grpc.Server, ip string) {
	var err error
	// 获取空闲端口
	freePort, err := utils.GetFreePort()
	if err != nil {
		zap.S().Errorw("get free port failed", "err", err.Error())
		return
	}
	global.Port = freePort
	zap.S().Infow("Info", "message", fmt.Sprintf("free port is : %d", global.Port))

	// proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, global.Port))
	if err != nil {
		zap.S().Errorw("net.Listen error", "err", err.Error())
		return
	}

	// 创建consul 客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)
	global.Client, err = api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create consul client failed", "err", err.Error())
		return
	}

	// 生成检查对象
	registerInfo := global.ServiceConfig.RegisterInfo
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServiceConfig.Host, global.Port),
		GRPCUseTLS:                     false,
		Timeout:                        registerInfo.CheckTimeOut,
		Interval:                       registerInfo.CheckInterval,
		DeregisterCriticalServiceAfter: registerInfo.DeregisterTime,
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServiceConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	global.ServiceID = serviceID
	registration.Tags = registerInfo.Tags
	registration.ID = serviceID
	registration.Port = global.Port
	registration.Address = global.ServiceConfig.Host
	registration.Check = check

	// 注册服务
	err = global.Client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("Error", "message", "client.Agent().ServiceRegister failed", "err", err.Error())
		return
	}
	zap.S().Infow("Info", "message", "service register success", "port", registration.Port, "ID", global.ServiceID)

	// 健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err = server.Serve(listen)
		panic(err)
	}()

	// 优雅停机
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 停机后注销服务
	err = global.Client.Agent().ServiceDeregister(global.ServiceID)
	if err != nil {
		zap.S().Errorw("global.Client.Agent().ServiceDeregister failed", "err", err.Error())
		return
	}
	zap.S().Infow("service deregister success", "serviceID", global.ServiceID)
}
