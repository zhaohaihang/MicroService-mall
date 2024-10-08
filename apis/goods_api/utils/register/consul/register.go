package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/zhaohaihang/goods_api/global"
)

type Registry struct {
	Host string
	Port int
}

type RegistryConfig struct {
	Address string
	Port    int
	Name    string
	Tags    []string
	Id      string
}

type RegistryClient interface {
	Register(config *RegistryConfig) error
	DeRegister(serviceId string) error
}

func NewRegistry(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(config *RegistryConfig) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ApiConfig.ConsulInfo.Host, global.ApiConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", config.Address, config.Port),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = config.Name
	registration.ID = config.Id
	registration.Port = config.Port
	registration.Tags = config.Tags
	registration.Address = config.Address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	client, err = api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
