package middleware

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type Times struct {
	TimeSlice int64 `json:"t_s"`
}

func TimeMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检测公共接口时间,时间差值5分钟
		// {debug}参数为720h
		Time := int64(300)
		if mod := func() string {
			var conf config.ConFig = &config.Config{}
			cond := conf.InitConfig()
			return cond.ConfData.Run.Mode
		}(); mod == "debug" {
			Time = int64(720 * 60 * 60)
		}
		json := Times{}
		_ = c.ShouldBindBodyWith(&json,binding.JSON)
		if time.Now().Unix() > json.TimeSlice + Time {
			c.JSON(errors.ErrTimeNoSwitch.HttpCode,gin.H{
				"code":errors.ErrTimeNoSwitch.Code,
				"message":errors.ErrTimeNoSwitch.Message,
			})
			c.Abort()
		} else {

		}
	}
}