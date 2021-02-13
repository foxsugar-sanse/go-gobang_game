package middleware

import (
	"fmt"
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

// jwt中间件认证
func JwtMiddlewareOAuth() gin.HandlerFunc {
	return func(c *gin.Context){
		var tk auth.JwtAPI = &auth.JWT{}
		tk.Init()
		headOauth :=c.Request.Header.Get("Authorization")
		tokenInfo := strings.SplitN(headOauth," ",2)
		// 验证信息为空则未登录
		if len(tokenInfo) != 2 && tokenInfo[0] != "Bearer" {
			c.JSON(errors.ErrTokenNotFound.HttpCode,gin.H{
				"code":errors.ErrTokenNotFound.Code,
				"message":errors.ErrTokenNotFound.Message,
			})
			c.Abort()
		}
		token,code := tk.MatchToken(tokenInfo[1])
		var md model.UserState = &model.OperationRedis{}
		bl := md.UserGetSignState(token.Uid)
		if code == true {
			if bl {
				goto DecodeTokenOk
			} else {
				c.JSON(errors.ErrUserSignNotFound.HttpCode,gin.H{
					"code":errors.ErrUserSignNotFound.Code,
					"message":errors.ErrUserSignNotFound.Message,
				})
				c.Abort()
			}
		} else if code == false && token == nil && bl == false{
			c.JSON(errors.ErrTokenValidation.HttpCode,gin.H{
				"code":errors.ErrTokenValidation.Code,
				"message":errors.ErrTokenValidation.Message,
			})
			//panic(errors.ErrTokenValidation.Message)
		}
		// 对比新时间与原始时间
		DecodeTokenOk:
		if time.Now().Unix() > token.ExpiresAt {
			// 获取redis中用户登录的状态并删除
			var d db.DB = &db.SetData{}
			dblink := d.RedisInit(1)
			_,err3 := dblink.Do("DEL",token.Uid).Result()
			if err3 == redis.Nil {
				// TODO : 未做日志
				fmt.Println("用户登录状态未存储或主动注销")
			}
			// 给前端返回信息
			c.JSON(errors.ErrTokenExpire.HttpCode,gin.H{
				"code":errors.ErrTokenExpire.Code,
				"message":errors.ErrTokenExpire.Message,
			})
		} else {

		}
	}

}

