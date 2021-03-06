package controller

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/foxsuagr-sanse/go-gobang_game/common/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type UserGroupController interface {
	UserGroupGet(c * gin.Context)
	UserGroupCreate(c * gin.Context)
	UserGroupSet(c * gin.Context)
	UserGroupDelete(c * gin.Context)
}

type OperationUserGroup struct {

}

type UserCreateGroupBindJSON struct {
	Group string `json:"group"`
}

type UserDeleteGroupBindJSON struct {
	Group string `json:"group"`
}

type UserSetGroupBindJSON struct {
	OldGroup string `json:"old_group"`
	NewGroup string `json:"new_group"`
}

func (o OperationUserGroup) UserGroupGet(c *gin.Context) {
	var md model.UserGroupDataBase = &model.OperationsUserGroup{}
	// 解密token获取用户Id
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if claims,bl := jwt.MatchToken(tokenInfo[1]); bl {
		uid,_ := strconv.ParseInt(claims.Uid,10,64)
		// 判断有无数据
		if groupSlices := md.GetUserGroup(uid); len(groupSlices) > 0 {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":groupSlices,
			})
		} else {
			c.JSON(errors.ErrGroupNotFound.HttpCode,gin.H{
				"code":errors.ErrGroupNotFound.Code,
				"message":errors.ErrGroupNotFound.Message,
			})
		}
	}
}

func (o OperationUserGroup) UserGroupCreate(c *gin.Context) {
	json := UserCreateGroupBindJSON{}
	_ = c.BindJSON(&json)
	var md model.UserGroupDataBase = &model.OperationsUserGroup{}
	// 解密token获取用户Id
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	utils.UserInput(json)
	if claims,bl := jwt.MatchToken(tokenInfo[1]); bl {
		uid,_ := strconv.ParseInt(claims.Uid,10,64)
		if md.AddUserGroup(uid,json.Group) {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
			})
		} else {
			c.JSON(errors.ErrGroupExist.HttpCode,gin.H{
				"code":errors.ErrGroupExist.Code,
				"message":errors.ErrGroupExist.Message,
			})
		}
	}
}

func (o OperationUserGroup) UserGroupSet(c *gin.Context) {
	json := UserSetGroupBindJSON{}
	_ = c.BindJSON(&json)
	// 设置分组名不能和默认分组名相同
	if (json.NewGroup != json.OldGroup) && json.OldGroup != model.DefaultUserGroupName && json.NewGroup != model.DefaultUserGroupName {
		var md model.UserGroupDataBase = &model.OperationsUserGroup{}
		// 解密token获取用户Id
		tokenHeader :=c.Request.Header.Get("Authorization")
		tokenInfo := strings.SplitN(tokenHeader, " ", 2)
		var jwt auth.JwtAPI = &auth.JWT{}
		jwt.Init()
		if claims,bl := jwt.MatchToken(tokenInfo[1]); bl {
			uid,_ := strconv.ParseInt(claims.Uid,10,64)
			if md.SetUserGroup(uid,json.OldGroup,json.NewGroup) {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
				})
			} else {
				c.JSON(errors.ErrGroupNotFound.HttpCode,gin.H{
					"code":errors.ErrGroupNotFound.Code,
					"message":errors.ErrGroupNotFound.Message,
				})
			}
		}
	} else {
		c.JSON(errors.ErrJsonArgFailed.HttpCode,gin.H{
			"code":errors.ErrJsonArgFailed.Code,
			"message":errors.ErrJsonArgFailed.Message,
		})
		c.Abort()
	}
}

func (o OperationUserGroup) UserGroupDelete(c *gin.Context) {
	//json := UserDeleteGroupBindJSON{}
	//_ = c.BindJSON(&json)
	group := c.Param("name")
	var md model.UserGroupDataBase = &model.OperationsUserGroup{}
	// 解密token获取用户Id
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if claims,bl := jwt.MatchToken(tokenInfo[1]); bl {
		uid,_ := strconv.ParseInt(claims.Uid,10,64)
		if md.DeleteUserGroup(uid,group) {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
			})
		} else {
			c.JSON(errors.ErrGroupNotFound.HttpCode,gin.H{
				"code":errors.ErrGroupNotFound.Code,
				"message":errors.ErrGroupNotFound.Message,
			})
		}
	}
}
