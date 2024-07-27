package initialize

import (
	"fmt"
	"github.com/zhaohaihang/goods_service/config"
	"github.com/zhaohaihang/goods_service/global"
	"path"
	"runtime"
)

// InitFileAbsPath 初始化文件路径
func InitFileAbsPath(cfgOption string) {
	basePath := getCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-"+cfgOption+".yaml",
		LogFile:    basePath + "/log",
	}
	fmt.Println("文件路径初始化成功:", basePath)
}

func getCurrentAbsolutePath() string {
	var abPath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(fileName)
	}
	return abPath
}
