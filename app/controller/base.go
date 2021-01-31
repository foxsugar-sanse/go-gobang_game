package controller

import (
	"github.com/gin-gonic/gin"
)

type RouterRequest interface {
	LoginPost(c * gin.Context)
	UserInfoGet(c * gin.Context)
	UserInfoUpdate(c * gin.Context)
	UserDelete(c * gin.Context)
	UserOtherOperations(c * gin.Context)
}

type UserRouter struct {
	
}

func (u UserRouter) LoginPost(c *gin.Context) {
	panic("implement me")
}

func (u UserRouter) UserInfoGet(c *gin.Context) {
	panic("implement me")
}

func (u UserRouter) UserInfoUpdate(c *gin.Context) {
	panic("implement me")
}

func (u UserRouter) UserDelete(c *gin.Context) {
	panic("implement me")
}

func (u UserRouter) UserOtherOperations(c *gin.Context) {
	panic("implement me")
}



