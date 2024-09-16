package initialize

import (
	"github.com/zhaohaihang/userop_api/config"
	"github.com/zhaohaihang/userop_api/global"
	"github.com/zhaohaihang/userop_api/utils"
)

func InitFileAbsPath(configFile string) {
	basePath := utils.GetCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-" + configFile + ".yaml",
		LogFile:    basePath + "/log",
	}
}
