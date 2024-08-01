package mode

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/zhaohaihang/userop_web/global"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
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
		HTTP:                           fmt.Sprintf("http://%s:%d/userop/v1/health", global.ApiConfig.Host, global.Port),
		GRPCUseTLS:                     false,
		Timeout:                        checkConfig.CheckTimeOut,
		Interval:                       checkConfig.CheckInterval,
		DeregisterCriticalServiceAfter: checkConfig.DeregisterTime,
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
		zap.S().Errorw("Error", "message", "userop_api register failed", "err", err.Error())
		return
	}
	zap.S().Infow("Info", "message", "userop_api register success", "port", registration.Port, "ID", ServiceID)

	client.Agent().AgentHealthServiceByID(serviceID)

	// 优雅停机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = client.Agent().ServiceDeregister(ServiceID)
	if err != nil {
		zap.S().Errorw("userop_api service deregister failed", "err", err.Error())
		return
	}
	zap.S().Infow("userop_api service deregister success", "serviceID", ServiceID)
}
