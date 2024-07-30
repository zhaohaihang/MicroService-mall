package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/hashicorp/consul/api"
	"github.com/zhaohaihang/inventory_service/config"
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
	Redsync       *redsync.Redsync
)
