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
	// Port := flag.Int("port", 8022, "服务启动端口")
	Mode := flag.String("mode", "release", "开发模式debug / 服务注册release")
	flag.Parse()
	initialize.InitFileAbsPath()
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitTranslator("zh")
	initialize.InitGoodsServiceConn()

	port, err := utils.GetFreePort()
	if err == nil {
		global.Port = port
	}

	go initialize.InitRouter()

	// 判断启动模式
	if *Mode == "debug" {
		zap.S().Warnf("debug本地调试模式 \n")
		mode.DebugMode()
	} else if *Mode == "release" {
		zap.S().Warnf("release服务注册模式 \n")
		mode.ReleaseMode()
	}
	
}
