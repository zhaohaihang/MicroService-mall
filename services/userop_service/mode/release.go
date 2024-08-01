package mode

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"github.com/zhaohaihang/userop_service/global"
	"github.com/zhaohaihang/userop_service/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"net"
	"os"
	"os/signal"
	"syscall"
)

func ReleaseMode(server *grpc.Server, ip string) {
	freePort, err := util.GetFreePort()
	if err != nil {
		zap.S().Errorw("get free port failed", "err", err.Error())
		return
	}
	global.FreePort = freePort
	zap.S().Infow("Info", "message", "free port is:", global.FreePort)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, global.FreePort))
	if err != nil {
		zap.S().Errorw("net.Listen failed", "err", err.Error())
		return
	}

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)
	global.Client, err = api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create new consul client failed", "err", err.Error())
		return
	}
	// 生成检查对象
	checkInfo := global.ServiceConfig.RegisterInfo
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServiceConfig.Host, global.FreePort),
		GRPCUseTLS:                     false,
		Timeout:                        checkInfo.CheckTimeOut,
		Interval:                       checkInfo.CheckInterval,
		DeregisterCriticalServiceAfter: checkInfo.DeregisterTime,
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServiceConfig.Name
	v4, err := uuid.NewV4()
	if err != nil {
		zap.S().Errorw("uuid.NewV4 failed", "err", err.Error())
		return
	}
	serviceID := v4.String()
	global.ServiceID = serviceID
	registration.ID = serviceID
	registration.Port = global.FreePort
	registration.Tags = checkInfo.Tags
	registration.Address = global.ServiceConfig.Host
	registration.Check = check
	err = global.Client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("userop_service register failed", "err", err.Error())
		return
	}
	zap.S().Infow("userop_service register success", "port", registration.Port, "ID", global.ServiceID)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	go func() {
		err = server.Serve(listen)
		panic(err)
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = global.Client.Agent().ServiceDeregister(global.ServiceID)
	if err != nil {
		zap.S().Errorw("userop_service deregister failed", "err", err.Error())
		return
	}
	zap.S().Infow("userop_service deregister success", "serviceID", global.ServiceID)
}
