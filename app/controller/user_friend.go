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

type CreateUserFriendBindJson struct {
	ReceiveId		int64 `json:"receive_id"`
	NoteInfo 		string `json:"note_info"`
}

type OperationUserFriendBindJson struct {
	Operation string `json:"op"`
	RepId	  int64  `json:"rep_id"`
}

type DeleteUserFriendBindJson struct {
	FriendUid int64 `json:"fid"`
}

type ModifyFriendInfoBindJson struct {
	FriendUid int64 `json:"fid"`
	Note	  string `json:"note"`
	Group	  string `json:"group"`
}

func (u *UserRouter) OtherUserFriendInterface(c *gin.Context) {

}

func (u *UserRouter) ModifyFriendInfo(c *gin.Context) {
	// 更改好友的信息（备注，分组）
	// 解密token
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	claims,_ := jwt.MatchToken(tokenInfo[1])
	json := ModifyFriendInfoBindJson{}
	_ = c.BindJSON(&json)
	var md model.User = &model.Operations{}
	uid,_ := strconv.ParseInt(claims.Uid,10,64)
	if md.SetUserFriendInfo(uid,json.FriendUid,json.Note,json.Group) {
		c.JSON(errors.OK.HttpCode,gin.H{
			"code":errors.OK.Code,
			"message":errors.OK.Message,
		})
	} else {
		c.JSON(errors.ErrUserFriendNotFound.HttpCode,gin.H{
			"code":errors.ErrUserFriendNotFound.Code,
			"message":errors.ErrUserFriendNotFound.Message,
		})
	}
}

func (u *UserRouter) DeleteUserForFriend(c *gin.Context) {
	// 删除好友
	// 解密token
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if claims,err := jwt.MatchToken(tokenInfo[1]);err {
		json := &DeleteUserFriendBindJson{}
		_ = c.BindJSON(&DeleteUserFriendBindJson{})
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(claims.Uid,10,64)
		if md.DeleteUserFriend(uid,json.FriendUid) {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
			})
		} else {
			c.JSON(errors.ErrUserFriendNotFound.HttpCode,gin.H{
				"code":errors.ErrUserFriendNotFound.Code,
				"message":errors.ErrUserFriendNotFound.Message,
			})
		}
	}
}

func (u *UserRouter) AddUserForFriend(c *gin.Context) {
	// TODO:废弃方法
}

func (u *UserRouter) GetUserForFriend(c *gin.Context) {
	// 获取该用户的所有好友
	// 解密token
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if claims,err := jwt.MatchToken(tokenInfo[1]);err {
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(claims.Uid,10,64)
		friendListMap := make(map[int]map[string]interface{})
		if friendList,bl := md.QueryUserFriend(uid); bl {
			for i := 0;i < len(friendList);i++ {
				// 组装数据
				friendListMap[i] = map[string]interface{}{
					"main_uid":friendList[i].MainUid,
					"friend_uid":friendList[i].FriendUid,
					"note":friendList[i].FriendNote,
					"group":friendList[i].UserGroup,
				}
			}
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":friendListMap,
			})
		} else {
			// 没有好友返回空map
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":friendListMap,
			})
		}
	}
}

func (u *UserRouter) RefuseUserFriendRequest(c *gin.Context) {
	// 拒绝好友申请
	// 获取jwt认证信息
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	// 解密token
	if claims,err := jwt.MatchToken(tokenInfo[1]);err {
		json := OperationUserFriendBindJson{}
		_ = c.BindJSON(&json)
		if json.Operation == "no" {
			var md model.UserForFriend = &model.OperationRedisForUf{}
			fid, _ := strconv.ParseInt(claims.Id,10,64)
			if md.UserFriendRequestRefuse(fid,json.RepId) {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
				})
			} else {
				c.JSON(errors.ErrUserFriendRequestFailed.HttpCode,gin.H{
					"code":errors.ErrUserFriendRequestFailed.Code,
					"message":errors.ErrUserFriendRequestFailed.Message,
				})
			}
		}
	}
}

func (u *UserRouter) ConsentUserFriendRequest(c *gin.Context) {
	// 获取jwt认证信息
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	// 解密token
	if claims,err := jwt.MatchToken(tokenInfo[1]);err {
		json := OperationUserFriendBindJson{}
		_ = c.BindJSON(&json)
		if json.Operation == "ok" {
			// 同意好友申请
			var md model.UserForFriend = &model.OperationRedisForUf{}
			fid, _ := strconv.ParseInt(claims.Uid,10,64)
			if md.UserFriendRequestConsent(fid,json.RepId) {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
				})
			}
		} else {
				c.JSON(errors.ErrUserFriendRequestFailed.HttpCode,gin.H{
					"code":errors.ErrUserFriendRequestFailed.Code,
					"message":errors.ErrUserFriendRequestFailed.Message,
				})
		}
	}
}

func (u *UserRouter) GetUserFriendRequest(c *gin.Context) {
	// 获取jwt认证信息
	tokenHeader :=c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	// 解密token
	if claims,err := jwt.MatchToken(tokenInfo[1]);err {
		var md model.UserForFriend = &model.OperationRedisForUf{}
		fid, _ := strconv.ParseInt(claims.Uid,10,64)
		// 获取好友申请
		if reFriendMap,err :=md.UserFriendRequestGet(fid);err {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":reFriendMap,
			})
		} else {
			c.JSON(errors.ErrUserFriendRequestFailed.HttpCode,gin.H{
				"code":errors.ErrUserFriendRequestFailed.Code,
				"message":errors.ErrUserFriendRequestFailed.Message,
			})
		}
	}
}

func (u *UserRouter) CreateUserFriendRequest(c *gin.Context) {
	json := CreateUserFriendBindJson{}
	_ = c.BindJSON(&json)
	// 解密token获得提交用户的uid
	tokenHeader := c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHeader, " ", 2)
	if tokenInfo[0] == "Bearer" && tokenInfo[1] != "" {
		var jwt auth.JwtAPI = &auth.JWT{}
		jwt.Init()
		if claims,bl := jwt.MatchToken(tokenInfo[1]);bl {
			var md model.UserForFriend = &model.OperationRedisForUf{}
			uid,_ := strconv.ParseInt(claims.Uid,10,64)
			if md.UserFriendRequestCreate(uid,json.ReceiveId,json.NoteInfo) {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
				})
			} else {
				c.JSON(errors.ErrUserFriendRequest.HttpCode,gin.H{
					"code":errors.ErrUserFriendRequest.Code,
					"message":errors.ErrUserFriendRequest.Message,
				})
			}
		}
	}

}

func (u *UserRouter) FormGroupGetUserFriend(c * gin.Context) {
	// 根据用户分组获取该用户的好友
	// TODO: 从Token中获取用户Uid的函数放在utils包里面了，简化代码
	groupName := c.Param("name")
	if claims,bl := utils.GinMatchToken(c);bl {
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(claims.Uid, 10,64)
		userFriendMap := make(map[int]int64)
		if UserFriendSlice,bl := md.FormGroupGetUserFriend(uid,groupName); bl {
			for i := 0;i < len(UserFriendSlice); i++ {
				userFriendMap[i + 1] = UserFriendSlice[i].FriendUid
			}
		}
		c.JSON(errors.OK.HttpCode,gin.H{
			"code":errors.OK.Code,
			"message":errors.OK.Message,
			"data":userFriendMap,
		})
	} else {
		c.JSON(errors.ErrGroupNotFound.HttpCode,gin.H{
			"code":errors.ErrGroupNotFound.Code,
			"message":errors.ErrGroupNotFound.Message,
		})
	}
}