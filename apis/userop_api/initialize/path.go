package initialize

import (
	"path"
	"runtime"

	"github.com/zhaohaihang/userop_api/config"
	"github.com/zhaohaihang/userop_api/global"
)

func InitFileAbsPath(configFile string) {
	basePath := getCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-" + configFile + ".yaml",
		LogFile:    basePath + "/log",
	}
}

func getCurrentAbsolutePath() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
