package controller

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type CreateUserFriendBindJson struct {
	SponSorId		int64 `json:"ss_id"`
	ReceiveId		int64 `json:"receive_id"`
	NoteInfo 		string `json:"note_info"`
}


func (u *UserRouter) OtherUserFriendInterface(c *gin.Context) {

}

func (u *UserRouter) ModifyFriendInfo(c *gin.Context) {

}

func (u *UserRouter) DeleteUserForFriend(c *gin.Context) {

}

func (u *UserRouter) AddUserForFriend(c *gin.Context) {

}

func (u *UserRouter) GetUserForFriend(c *gin.Context) {

}

func (u *UserRouter) RefuseUserFriendRequest(c *gin.Context) {

}

func (u *UserRouter) ConsentUserFriendRequest(c *gin.Context) {

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
		fid, _ := strconv.ParseInt(claims.Id,10,64)
		// 获取好友申请
		if reFriendMap,err :=md.UserFriendRequestGet(fid);err {
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":map[string]string{
					"main_uid":reFriendMap["main_uid"],
					"friend_uid":reFriendMap["friend_uid"],
					"note":reFriendMap["note"],
					"state":reFriendMap["state"],
				},
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
	json := &CreateUserFriendBindJson{}
	_ = c.BindJSON(&CreateUserFriendBindJson{})
	var md model.UserForFriend = &model.OperationRedisForUf{}
	if md.UserFriendRequestCreate(json.SponSorId,json.ReceiveId,json.NoteInfo) {
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