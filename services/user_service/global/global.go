package global

import (
	"github.com/hashicorp/consul/api"
	"github.com/zhaohaihang/user_service/config"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	ServiceConfig *config.ServiceConfig
	FilePath      *config.FilePathConfig
	NacosConfig   *config.NacosConfig
	Port          int
	Client        *api.Client
	ServiceID     string
)
