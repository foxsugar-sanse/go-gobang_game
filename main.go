package main

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/service"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/utils"
	"github.com/foxsuagr-sanse/go-gobang_game/router"
	"github.com/foxsuagr-sanse/go-gobang_game/router/middleware"
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
	r.Use(middleware.Cors())
	// 初始化路由
	r.MaxMultipartMemory = 64 << 20 // 64MiB 上传时用到的内存大小
	var x router.Router = &router.Route{}
	// 创建连接管道
	var wsChan1 = make(chan *service.ClientMessage,100)
	var wsChan2 = make(chan *service.ClientMessage,100)

	// 通信管道
	var wsChan3 = make(chan string)
	var wsChan4 = make(chan string)

	x.Run(r,wsChan1,wsChan2,wsChan3,wsChan4)
	go func() {
		// 维护webSocket连接
		for  {
			select {
			case cm := <- wsChan1:
				select {
				case cm2 := <- wsChan2:
					if cm2.OldUid == cm.OldUid {
						var sr service.Service = &service.Server{}
						go sr.GameService(cm,cm2,wsChan3,wsChan4)
					} else {
						// 匹配不成功把数据放回管道
						wsChan1 <- cm
						wsChan2 <- cm2
					}
				}

			}
		}
	}()
	_ = r.Run(cond.ConfData.Run.Ipaddr + cond.ConfData.Run.Port)
}
