package mode

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"

	"github.com/zhaohaihang/inventory_service/global"
	"github.com/zhaohaihang/inventory_service/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func ReleaseMode(server *grpc.Server, ip string) {
	freePort, err := util.GetFreePort()
	if err != nil {
		zap.S().Errorw("get free port failed", "err", err.Error())
		return
	}
	global.FreePort = freePort
	zap.S().Infow("get free port success", "message", "free port is:", global.FreePort)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, global.FreePort))
	if err != nil {
		zap.S().Errorw("net.Listen failed", "err", err.Error())
		return
	}

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServiceConfig.ConsulInfo.Host, global.ServiceConfig.ConsulInfo.Port)
	global.Client, err = api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create new consul client failed", "err", err.Error())
		return
	}
	checkInfo := global.ServiceConfig.RegisterInfo
	// 生成检查对象
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
	serviceID := uuid.NewV4().String()
	global.ServiceID = serviceID
	registration.ID = serviceID
	registration.Port = global.FreePort
	registration.Tags = checkInfo.Tags
	registration.Address = global.ServiceConfig.Host
	registration.Check = check
	err = global.Client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("Error", "message", "inventory_service register failed", "err", err.Error())
		return
	}
	zap.S().Infow("Info", "message", "inventory_service register success", "port", registration.Port, "ID", global.ServiceID)


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
		zap.S().Errorw("inventory_service deregister failed", "err", err.Error())
		return
	}
	zap.S().Infow("inventory_service deregister success", "serviceID", global.ServiceID)
}
