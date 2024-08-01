package main

import (
	"flag"

	"github.com/zhaohaihang/userop_web/global"
	"github.com/zhaohaihang/userop_web/initialize"
	"github.com/zhaohaihang/userop_web/mode"
	"github.com/zhaohaihang/userop_web/utils"
	"go.uber.org/zap"
)

func main() {
	Mode := flag.String("mode", "release", "mode debug / release ")
	flag.Parse()

	initialize.InitFileAbsPath()
	initialize.InitConfig()
	initialize.InitLogger()
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
