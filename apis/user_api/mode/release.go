package mode

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"github.com/zhaohaihang/user_api/global"
	"go.uber.org/zap"
)

func ReleaseMode() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ApiConfig.ConsulInfo.Host, global.ApiConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorw("create new consul client failed", "err", err.Error())
		return
	}

	checkConfig := global.ApiConfig.ServiceInfo
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ApiConfig.Host, global.Port),
		GRPCUseTLS:                     false,
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ApiConfig.Name
	serviceID := uuid.NewV4().String()
	ServiceID := serviceID
	registration.ID = serviceID
	registration.Port = global.Port
	registration.Tags = checkConfig.Tags
	registration.Address = global.ApiConfig.Host
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorw("Error", "message", "user_api register failed", "err", err.Error())
		return
	}
	zap.S().Infow("Info", "message", "user_api register success", "port", registration.Port, "ID", ServiceID)

	client.Agent().AgentHealthServiceByID(serviceID)

	// 优雅停机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("serviceID", ServiceID)
	err = client.Agent().ServiceDeregister(ServiceID)
	if err != nil {
		zap.S().Errorw("user_api service deregister  failed", "err", err.Error())
		return
	}
	zap.S().Infow("user_api service deregister success", "serviceID", ServiceID)
}
