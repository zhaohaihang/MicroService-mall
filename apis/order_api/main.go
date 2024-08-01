package main

import (
	"flag"

	"github.com/zhaohaihang/order_api/global"
	"github.com/zhaohaihang/order_api/initialize"
	"github.com/zhaohaihang/order_api/mode"
	"github.com/zhaohaihang/order_api/utils"
	"go.uber.org/zap"
)

func main() {
	Mode := flag.String("mode", "release", "mode debug / release ")
	flag.Parse()
	// 初始化文件路径
	initialize.InitFileAbsPath()
	initialize.InitLogger()
	initialize.InitConfig()
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
