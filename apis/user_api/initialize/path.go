package initialize

import (
	"github.com/zhaohaihang/user_api/config"
	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/utils"
)

// InitFilePath 初始化全局文件路径
func InitFilePath(nacosConfig string) {
	basePath := utils.GetCurrentAbsolutePath()
	global.FileConfig = &config.FileConfig{
		ConfigFile: basePath + "/config-" + nacosConfig + ".yaml",
		LogFile:    basePath + "/log",
	}
}
