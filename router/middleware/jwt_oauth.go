package middleware

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// jwt中间件认证
func JwtMiddlewareOAuth() gin.HandlerFunc {
	var tk auth.JwtAPI = &auth.JWT{}
	tk.Init()
	return func(c *gin.Context){
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
		if code == true {
			goto DecodeTokenOk
		} else if code == false && token == nil{
			c.JSON(errors.ErrTokenValidation.HttpCode,gin.H{
				"code":errors.ErrTokenValidation.Code,
				"message":errors.ErrTokenValidation.Message,
			})
			c.Abort()
			panic(errors.ErrTokenValidation.Message)
		}
		// 对比新时间与原始时间
		DecodeTokenOk:
		if time.Now().Unix() > token.ExpiresAt {
			c.JSON(errors.ErrTokenExpire.HttpCode,gin.H{
				"code":errors.ErrTokenExpire.Code,
				"message":errors.ErrTokenExpire.Message,
			})
		} else {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
			})
		}
	}

}

