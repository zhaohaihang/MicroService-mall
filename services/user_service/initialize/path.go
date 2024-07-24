package initialize

import (
	"fmt"
	"path"
	"runtime"

	"github.com/zhaohaihang/user_service/config"
	"github.com/zhaohaihang/user_service/global"
)

// InitFileAbsPath 初始化文件路径
func InitFileAbsPath() {
	basePath := getCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basePath + "/config-debug.yaml",
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
