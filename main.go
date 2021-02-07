package main

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	r := gin.New()
	// 初始化路由
	var x router.Router = &router.Route{}
	x.Run(r)

	_ = r.Run(cond.ConfData.Run.Ipaddr + cond.ConfData.Run.Port)
}
