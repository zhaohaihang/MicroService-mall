package main

import (
	"flag"

	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/initialize"
	"github.com/zhaohaihang/user_api/mode"
	"github.com/zhaohaihang/user_api/utils"
	"go.uber.org/zap"
)

func main() {
	Mode := flag.String("mode", "release", "mode debug / release ")
	flag.Parse()
	initialize.InitFilePath()
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitTranslator("zh")
	initialize.InitValidator()
    initialize.InitUserService()

	port, err := utils.GetFreePort()
	if err == nil {
		global.Port = port
	}
	go initialize.InitRouters()

	if *Mode == "release" {
		zap.S().Warnf("release服务注册模式 \n")
		mode.ReleaseMode()
	} else if *Mode == "debug" {
		zap.S().Warnf("debug本地调试模式 \n")
		mode.DebugMode()
	}

}
