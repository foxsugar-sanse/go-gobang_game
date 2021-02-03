package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
)

type Times struct {
	TimeSlice int64 `json:"t_s"`
}

func TimeMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检测公共接口时间,时间差值5分钟
		json := &Times{}
		_ = c.BindJSON(&json)
		if time.Now().Unix() > json.TimeSlice + 300 {
			c.JSON(errors.ErrTimeNoSwitch.HttpCode,gin.H{
				"code":errors.ErrTimeNoSwitch.Code,
				"message":errors.ErrTimeNoSwitch.Message,
			})
			c.Abort()
		}
	}
}