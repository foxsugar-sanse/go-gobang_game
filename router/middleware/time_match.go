package middleware

import (
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
		json := Times{}
		_ = c.ShouldBindBodyWith(&json,binding.JSON)
		if time.Now().Unix() > json.TimeSlice + 300 {
			c.JSON(errors.ErrTimeNoSwitch.HttpCode,gin.H{
				"code":errors.ErrTimeNoSwitch.Code,
				"message":errors.ErrTimeNoSwitch.Message,
			})
			c.Abort()
		} else {
			c.Next()
		}
	}
}