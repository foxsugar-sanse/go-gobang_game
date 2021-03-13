package controller

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/app/service"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/foxsuagr-sanse/go-gobang_game/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ymzuiku/hit"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserGames interface {
	CreateUserInvite	   (c * gin.Context)
	RefuseUserInvite	   (c * gin.Context)
	GetUserInvite		   (c * gin.Context)
	ConSentInvite		   (c * gin.Context)
	DeleteUserInvite	   (c * gin.Context)
	ResponseWebSocketLink  (c * gin.Context,wsChan1 chan *service.ClientMessage,wsChan2 chan *service.ClientMessage,wsChan3 chan string,wsChan4 chan string)
	MainUidCreateWebSocket (c * gin.Context)
	AceUidCreateWebSocket  (c * gin.Context)
}

type UserInviteFunc struct {

}

type InviteMessageStrut struct {
	MessageID     int64		`jsons:"MessageId"`
	MainWsId      int64		`jsons:"MainWsId"`
	AceWsId       int64		`json:"AceWsId"`
	OldUid 		  int64 	`json:"OldUid"`
	AceUid        int64 	`json:"AceUid"`
	Info          string 	`json:"Info"`
	MessageStatus string	`json:"MessageStatus"`
	MessageState  int		`json:"MessageState"` // 0为创建时,1为同意,2为拒绝
}

type InviteBindJson struct {
	OldUid 		  int64 	`json:"main_uid"`
	AceUid        int64 	`json:"accept_uid"`
	Info          string 	`json:"invite_info"`
}

func (u UserInviteFunc) CreateUserInvite(c *gin.Context) {
	// TODO:多个游戏邀请信息
	jsons := InviteBindJson{}
	_ = c.BindJSON(&jsons)
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if token, bl := jwt.MatchToken(tokenHeadInfo[1]); bl {
		// 填充结构体中的信息
		// 初始化结构体
		invStru := InviteMessageStrut{
			MessageID:     0,
			MainWsId:      0,
			AceWsId:       0,
			OldUid:        jsons.OldUid,
			AceUid:        jsons.AceUid,
			Info:          jsons.Info,
			MessageStatus: "",
			MessageState:  0,
		}
		oldUid, _ := strconv.ParseInt(token.Uid,10,64)
		invStru.MainWsId = jsons.OldUid
		invStru.OldUid = oldUid
		invStru.MessageStatus = "start"
		invStru.MessageState = 0
		// 使用雪花算法生成消息id
		sf,_ := snowflake.NewNode(0)
		invStru.MessageID = int64(sf.Generate())
		res,_ := json.Marshal(invStru)
		var md model.InviteRedis = &model.Invite{}

		errno := hit.If(md.CreateInvite(strconv.FormatInt(jsons.AceUid, 10), string(res)),errors.OK,errors.ErrRedis).(*errors.Errno)
		c.JSON(errno.HttpCode,gin.H{
			"code":errno.Code,
			"message":errno.Message,
		})
	}
}

func (u UserInviteFunc) RefuseUserInvite(c *gin.Context) {
	// 拒绝游戏申请
	OldUid := c.Query("aid")
	var md model.InviteRedis = &model.Invite{}
	claims,bl := utils.GinMatchToken(c)
	mid,_ := strconv.ParseInt(claims.Uid,10,64)
	err := hit.If(bl,func() *errors.Errno {
		bl := md.DeleteInvite(OldUid,mid)
		return hit.If(bl,errors.OK,errors.ErrUserInviteNotFound).(*errors.Errno)
	},func () *errors.Errno {
		return errors.ErrUserSignNotFound
	})
	errno := err.(*errors.Errno)
	c.JSON(errno.HttpCode,gin.H{
		"code":errno.Code,
		"message":errno.Message,
	})
}

func (u UserInviteFunc) GetUserInvite(c *gin.Context) {
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead, " ", 2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	if token, bl := jwt.MatchToken(tokenHeadInfo[1]); bl {
		var md model.InviteRedis = &model.Invite{}
		var mp = make(map[int]InviteMessageStrut)
		var num int = 1
		if uidSlices,bl := md.GetInvite(token.Uid) ; bl {
			for i := 0;i < len(uidSlices);i++ {
				if uidSlices[i] != "" {
					invStrut := InviteMessageStrut{}
					if err := json.Unmarshal([]byte(uidSlices[i]),&invStrut) ; err == nil {
						mp[num] = invStrut
						num++
					} else {

					}
				}
			}
			// 包装完数据就返回
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":mp,
			})
		} else {
			c.JSON(errors.ErrUserInviteNot.HttpCode,gin.H{
				"code":errors.ErrUserInviteNot.Code,
				"message":fmt.Sprintf("{%s}",token.Uid) + errors.ErrUserInviteNot.Message,
			})
		}
	}
}

func (u UserInviteFunc) ConSentInvite(c *gin.Context) {
	panic("implement me")
}

func (u UserInviteFunc) DeleteUserInvite(c *gin.Context) {
	panic("implement me")
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r * http.Request) bool {
		return true
	},
}

func (u UserInviteFunc) ResponseWebSocketLink(c *gin.Context,wsChan1 chan *service.ClientMessage,wsChan2 chan *service.ClientMessage,wsChan3 chan string,wsChan4 chan string) {
	// 判断Query参数并接收对应的websocket连接
	op := c.Query("op")
	//var md model.InviteRedis = &model.Invite{}
	claims,bl := utils.GinMatchToken(c)
	if bl == false {
		log.Println("用户Token错误，禁止连接")
		panic("用户Token错误，禁止连接")
	}
	switch op {
	case "main":
		ws,err := upGrader.Upgrade(c.Writer,c.Request,nil)
		ws.Close() // TODO:站位
		if err != nil {
			return
		}
		//var wsService service.Service = &service.Server{}
		//go wsService.GoMainWebSocketService(ws,claims)
		ClientMessage := service.ClientMessage{
			OldUid:    claims.Uid,
			State:     0,
		}
		wsChan1 <- &ClientMessage
		if <- wsChan3 == "ok" {
			// 接收ok开始服务
		}
		break
	case "ace":
		// redis查询邀请 -> websocket通过
		var md model.InviteRedis = &model.Invite{}
		invSlice,_ := md.GetInvite(claims.Uid)
		for i := 0 ; i < len(invSlice); i++ {
			if invSlice[i] != "" {
				// 绑定Json
				jsons := InviteMessageStrut{}
				_ = json.Unmarshal([]byte(invSlice[i]), &jsons)
				// 邀请olduid和aceuid分别匹配时，状态为1则代表同意邀请
				if jsons.MessageState == 1 && strconv.FormatInt(jsons.AceUid,10) == claims.Uid {
					ws,err := upGrader.Upgrade(c.Writer,c.Request,nil)
					ws.Close() // TODO:占位
					if err != nil {
						return
					}
					wsChan2 <- &service.ClientMessage{
						MessageId: jsons.MessageID,
						OldUid:    strconv.FormatInt(jsons.OldUid,10),
						State:     jsons.MessageState,
						AceUid:    strconv.FormatInt(jsons.AceUid,10),
					}
					if <- wsChan4 == "ok" {
						// 接收ok开始服务
					}
				}
			}
		}
		break
	default:
		// 参数错误

	}
}

func (u UserInviteFunc) MainUidCreateWebSocket(c *gin.Context) {
	panic("implement me")
}

func (u UserInviteFunc) AceUidCreateWebSocket(c *gin.Context) {
	panic("implement me")
}
