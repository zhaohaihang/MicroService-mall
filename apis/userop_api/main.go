package main

import (
	"flag"

	"github.com/zhaohaihang/userop_api/global"
	"github.com/zhaohaihang/userop_api/initialize"
	"github.com/zhaohaihang/userop_api/mode"
	"github.com/zhaohaihang/userop_api/utils"
	"go.uber.org/zap"
)

func main() {
	Mode := flag.String("mode", "release", "mode debug / release ")
	nacosConfig := flag.String("nacosConfig", "home", "home / com ")
	flag.Parse()

	initialize.InitFileAbsPath(*nacosConfig)
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSentinel()
	initialize.InitTranslator("zh")
	initialize.InitValidator()
	initialize.InitService()

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
