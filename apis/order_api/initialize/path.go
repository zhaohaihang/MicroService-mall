package initialize

import (
	"fmt"

	"github.com/zhaohaihang/order_api/config"
	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/utils"
)

func InitFileAbsPath(configFile string) {
	basePath := utils.GetCurrentAbsolutePath()
	global.FilePathConfig = &config.FilePathConfig{
		ConfigFile: basePath + "/config-" + configFile + ".yaml",
		LogFile:    basePath + "/log",
	}
	fmt.Println("init success:", basePath)
}
