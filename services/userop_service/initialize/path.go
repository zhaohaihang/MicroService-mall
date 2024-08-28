package initialize

import (
	"fmt"
	"path"
	"runtime"
	"github.com/zhaohaihang/userop_service/config"
	"github.com/zhaohaihang/userop_service/global"
)

func InitFileAbsPath(cfgOption string) {
	basePath := getCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-"+cfgOption+".yaml",
		LogFile:    basePath + "/log",
	}
	fmt.Println("file path init success:", basePath)
}

func getCurrentAbsolutePath() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
