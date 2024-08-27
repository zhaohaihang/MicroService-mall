package initialize

import (
	"fmt"
	"path"
	"runtime"

	"github.com/zhaohaihang/order_api/config"
	"github.com/zhaohaihang/order_api/global"
)

func InitFileAbsPath(configFile string) {
	basePath := getCurrentAbsolutePath()
	global.FilePathConfig = &config.FilePathConfig{
		ConfigFile: basePath + "/config-" + configFile + ".yaml",
		LogFile:    basePath + "/log",
	}
	fmt.Println("init success:", basePath)
}

func getCurrentAbsolutePath() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
