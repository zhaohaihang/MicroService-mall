package initialize

import (
	"fmt"
	"github.com/zhaohaihang/inventory_service/config"
	"github.com/zhaohaihang/inventory_service/global"
	"path"
	"runtime"
)

func InitFileAbsPath() {
	basePath := getCurrentAbsplutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-debug.yaml",
		LogFile:    basePath + "/log",
	}
	fmt.Println("file path init success:", basePath)
}

func getCurrentAbsplutePath() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
