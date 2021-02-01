package controller

import (
	"github.com/gin-gonic/gin"
)

type RouterRequest interface {
	LoginPost			(c * gin.Context)
	UserInfoGet			(c * gin.Context)
	UserInfoUpdate		(c * gin.Context)
	UserDelete			(c * gin.Context)
	UserOtherOperations	(c * gin.Context)
}

type UserJsonBindLogin struct {
	TimeSlices 		int64  `json:"t_s"`
	UserName 		string `json:"user_name"`
	UserPassword 	string `json:"user_password"`
}

type UserRouter struct {
	UserJBL *UserJsonBindLogin
}

func (u *UserRouter) LoginPost(c *gin.Context) {
	// 鉴别登录状态


	_ = c.BindJSON(&UserJsonBindLogin{})
	panic("implement me")
}

func (u *UserRouter) UserInfoGet(c *gin.Context) {
	c.Param("name")
}

func (u *UserRouter) UserInfoUpdate(c *gin.Context) {
	panic("implement me")
}

func (u *UserRouter) UserDelete(c *gin.Context) {
	panic("implement me")
}

func (u *UserRouter) UserOtherOperations(c *gin.Context) {
	panic("implement me")
}



