package utils

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/gin-gonic/gin"
	"strings"
)

func GinMatchToken(c * gin.Context) (*auth.MyClaims,bool) {
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if claims,bl := jwt.MatchToken(tokenInfo[1]);bl {
		return claims,bl
	} else {
		return nil, false
	}
}