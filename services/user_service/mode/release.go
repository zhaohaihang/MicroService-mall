package mode

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"github.com/zhaohaihang/user_service/global"
	"github.com/zhaohaihang/user_service/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// ReleaseMode release服务注册模式
func ReleaseMode(server *grpc.Server, ip string) {
	freePort, err := util.GetFreePort()
	if err != nil {
		zap.S().Errorw("get free port failed", "err", err.Error())
		return
	}
	global.Port = freePort
	zap.S().Infow("Info", "message", "free port is:", global.Port)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, global.Port))
	if err != nil {
		zap.S().Errorw("net.Listen failed", "err", err.Error())
		return
	}

	// 生成检查对象
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)
	global.Client, err = api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create new consul client failed", "err", err.Error())
		return
	}
	// 生成检查对象
	checkConfig := global.ServiceConfig.ServiceInfo
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServiceConfig.Host, global.Port),
		GRPCUseTLS:                     false,
		Timeout:                        checkConfig.CheckTimeOut,
		Interval:                       checkConfig.CheckInterval,
		DeregisterCriticalServiceAfter: checkConfig.DeregisterTime,
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServiceConfig.Name
	serviceID := uuid.NewV4().String()
	global.ServiceID = serviceID
	registration.ID = serviceID
	registration.Port = global.Port
	registration.Tags = checkConfig.Tags
	registration.Address = global.ServiceConfig.Host
	registration.Check = check
	err = global.Client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("Error", "message", "user_service register failed", "err", err.Error())
		return
	}
	zap.S().Infow("Info", "message", "user_service register success", "port", registration.Port, "ID", global.ServiceID)

	// 健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err = server.Serve(listen)
		panic(err)
	}()

	// 优雅停机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = global.Client.Agent().ServiceDeregister(global.ServiceID)
	if err != nil {
		zap.S().Errorw("user_service deregister failed", "err", err.Error())
		return
	}
	zap.S().Infow("user_service deregister success", "serviceID", global.ServiceID)
}
