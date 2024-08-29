package main

import (
	"flag"
	"github.com/zhaohaihang/goods_api/global"
	"github.com/zhaohaihang/goods_api/initialize"
	"github.com/zhaohaihang/goods_api/mode"
	"github.com/zhaohaihang/goods_api/utils"

	"go.uber.org/zap"
)

func main() {
	Mode := flag.String("mode", "release", "mode debug / release ")
	nacosConfig := flag.String("nacosConfig","home","home / com")
	flag.Parse()

	initialize.InitFileAbsPath(*nacosConfig)
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitSentinel()
	initialize.InitTranslator("zh")
	initialize.InitGoodsServiceConn()

	port, err := utils.GetFreePort()
	if err == nil {
		global.Port = port
	}

	go initialize.InitRouter()
	
	// 判断启动模式
	if *Mode == "release" {
		zap.S().Warnf("release mode \n")
		mode.ReleaseMode()
	} else if *Mode == "debug" {
		zap.S().Warnf("debug mode \n")
		mode.DebugMode()
	}
	
}
