package global

import (
	"github.com/hashicorp/consul/api"
	"github.com/zhaohaihang/userop_service/config"
	"gorm.io/gorm"
)

var (
	FilePath      *config.FilePathConfig
	NacosConfig   *config.NacosConfig
	ServiceConfig *config.ServiceConfig
	DB            *gorm.DB
	Client        *api.Client
	FreePort      int
	ServiceID     string
)
