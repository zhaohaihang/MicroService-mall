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
	nacosConfig := flag.String("nacosConfig","home","home / com")
	flag.Parse()

	initialize.InitFilePath(*nacosConfig)
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitTranslator("zh")
	initialize.InitValidator()
	initialize.InitRedis()
    initialize.InitUserService()

	port, err := utils.GetFreePort()
	if err == nil {
		global.Port = port
	}
	go initialize.InitRouters()

	if *Mode == "release" {
		zap.S().Warnf("release mode \n")
		mode.ReleaseMode()
	} else if *Mode == "debug" {
		zap.S().Warnf("debug mode \n")
		mode.DebugMode()
	}

}
