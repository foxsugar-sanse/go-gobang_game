package service

import (
	"encoding/json"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

type Service interface {
	GoMainWebSocketService(conn *websocket.Conn,claims *auth.MyClaims) *errors.Errno
	GoAceWebSocketService(conn *websocket.Conn,claims *auth.MyClaims) *errors.Errno
	GameService(cm *ClientMessage,cm2 *ClientMessage,wsChan3 chan string,wsChan4 chan string) *errors.Errno
}

type Server struct {

}

type ClientMessage struct {
	MessageId     int64  // 对局id
	OldUid        string // 发起邀请的用户
	State         int    // 对局状态
	AceUid        string // 接受邀请的用户

	board         [19][19]int // 棋盘
	opening       string // 谁开局
	totalNumber   int    // 总步数
	totalTimer    int64    // 总时间，按照秒速记时
	blackChessman string // 黑棋
	whileChessman string // 白棋
	oldUidRecord  map[int]string // 发起邀请的用户记录
	aceUidRecord  map[int]string // 接收邀请的用户记录
}

type WebSocketClientMessage struct {
	Uid string 		`json:"uid"`
	Message string  `json:"message"`
}

type WebSocketServerMessage struct {
	MessageId 		int64
	Board 			[19][19]int
	BlackChessman 	string
	WhiteChessman 	string
	Info 			string
	Operation   	string
	Opening            string
	TotalNumber   	int
}

func (s Server) GoMainWebSocketService(conn *websocket.Conn, claims *auth.MyClaims) *errors.Errno {
	return errors.OK
}

func (s Server) GoAceWebSocketService(conn *websocket.Conn, claims *auth.MyClaims) *errors.Errno {
	panic("implement me")
}

func (s Server) GameService(cm *ClientMessage,cm2 *ClientMessage,wsChan3 chan string,wsChan4 chan string) *errors.Errno {
	// 数据初始化
	clientMessage := ClientMessage{
		MessageId: cm2.MessageId,
		OldUid:    cm.OldUid,
		State:     cm2.State,
		AceUid:    cm2.AceUid,
	}
	clientMessage.board = [19][19]int{}
	clientMessage.totalNumber = 0
	clientMessage.totalTimer = int64(time.Minute * 0)
	clientMessage.oldUidRecord = make(map[int]string)
	clientMessage.aceUidRecord = make(map[int]string)
	// 随机数选择开局者
	if num := rand.Intn(11) ; num / 2 > 5 {
		clientMessage.opening = cm2.AceUid
		clientMessage.whileChessman = cm2.AceUid
		clientMessage.blackChessman = cm2.OldUid
	} else {
		clientMessage.opening = cm2.OldUid
		clientMessage.whileChessman = cm2.OldUid
		clientMessage.blackChessman = cm2.AceUid
	}
	wsChan3 <- "ok"
	wsChan4 <- "ok"
	for {
		wsServerMessage := WebSocketServerMessage{
			MessageId:     clientMessage.MessageId,
			Board:         clientMessage.board,
			BlackChessman: clientMessage.blackChessman,
			WhiteChessman: clientMessage.whileChessman,
		}
		wsServerMessage.Opening = clientMessage.opening
		bytes,_ := json.Marshal(&wsServerMessage)
		wsChan3 <- string(bytes)
		wsChan4 <- string(bytes)

	}
}
