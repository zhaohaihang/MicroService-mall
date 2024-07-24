package main

import (
	"flag"

	"github.com/zhaohaihang/user_api/global"
	"github.com/zhaohaihang/user_api/initialize"
	"github.com/zhaohaihang/user_api/mode"
	"github.com/zhaohaihang/user_api/utils"
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
		mode.ReleaseMode()
	} else if *Mode == "debug" {
		mode.DebugMode()
	}

}
