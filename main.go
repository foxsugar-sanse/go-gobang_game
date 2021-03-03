package main

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/utils"
	"github.com/foxsuagr-sanse/go-gobang_game/router"
	"github.com/gin-gonic/gin"
)

func main() {

	// 读取配置文件
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	if cond.ConfData.Run.Mode == "debug" {
		go utils.ONEXIT()
	}
	r := gin.Default()
	// 初始化路由
	r.MaxMultipartMemory = 64 << 20 // 64MiB 上传时用到的内存大小
	var x router.Router = &router.Route{}
	x.Run(r)

	_ = r.Run(cond.ConfData.Run.Ipaddr + cond.ConfData.Run.Port)
}
