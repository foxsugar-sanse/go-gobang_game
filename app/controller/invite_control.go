package controller

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/ymzuiku/hit"
	"strconv"
	"strings"
)

type UserGames interface {
	CreateUserInvite	(c * gin.Context)
	RefuseUserInvite	(c * gin.Context)
	GetUserInvite		(c * gin.Context)
	ConSentInvite		(c * gin.Context)
	DeleteUserInvite	(c * gin.Context)
}

type UserInviteFunc struct {

}

type InviteMessageStrut struct {
	MessageID     int64
	MainWsId      int64
	AceWsId       int64
	OldUid 		  int64 	`json:"main_uid"`
	AceUid        int64 	`json:"accept_uid"`
	Info          string 	`json:"invite_info"`
	MessageStatus string
	MessageState  int
}

func (u UserInviteFunc) CreateUserInvite(c *gin.Context) {
	// TODO:多个游戏邀请信息
	jsons := InviteMessageStrut{}
	_ = c.BindJSON(&jsons)
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if token, bl := jwt.MatchToken(tokenHeadInfo[1]); bl {
		// 填充结构体中的信息
		oldUid, _ := strconv.ParseInt(token.Uid,10,64)
		jsons.MainWsId = jsons.OldUid
		jsons.OldUid = oldUid
		jsons.MessageStatus = "start"
		jsons.MessageState = 0
		// 使用雪花算法生成id
		sf,_ := snowflake.NewNode(0)
		jsons.MessageID = int64(sf.Generate())
		res,_ := json.Marshal(jsons)
		var md model.InviteRedis = &model.Invite{}

		errno := hit.If(md.CreateInvite(strconv.FormatInt(jsons.AceUid, 10),hex.EncodeToString(res)) == true,errors.OK,errors.ErrRedis).(*errors.Errno)
		c.JSON(errno.HttpCode,gin.H{
			"code":errno.Code,
			"message":errno.Message,
		})
	}
}

func (u UserInviteFunc) RefuseUserInvite(c *gin.Context) {
	panic("implement me")
}

func (u UserInviteFunc) GetUserInvite(c *gin.Context) {

}

func (u UserInviteFunc) ConSentInvite(c *gin.Context) {
	panic("implement me")
}

func (u UserInviteFunc) DeleteUserInvite(c *gin.Context) {
	panic("implement me")
}
