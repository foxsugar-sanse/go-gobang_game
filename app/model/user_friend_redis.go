package model

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"strconv"
)


type UserForFriend interface {
	UserFriendRequestCreate(mainUid int64,friendUid int64,note string) 	bool
	UserFriendRequestGet(friendUid int64) 								(map[string]string,bool)
	UserFriendRequestConsent(friendUid int64)
	UserFriendRequestRefuse(friendUid int64)
}

type OperationRedisForUf struct {

}

func (o OperationRedisForUf) UserFriendRequestCreate(mainUid int64, friendUid int64, note string) bool{
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(2)
	defer dblink.Close()
	// key格式"friendUid"
	err := dblink.HMSet(strconv.FormatInt(friendUid, 10), map[string]interface{}{
		"main_uid": mainUid,
		"friend_uid": friendUid,
		"note": note,
		"state":0,
	}).Err()
	return err != nil
}

func (o OperationRedisForUf) UserFriendRequestGet(friendUid int64) (map[string]string,bool) {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(2)
	defer dblink.Close()
	if requestMap, err := dblink.HGetAll(strconv.FormatInt(friendUid, 10)).Result();err == nil {
		return requestMap,true
	} else {
		return nil,false
	}
}

func (o OperationRedisForUf) UserFriendRequestConsent(friendUid int64) {
	// 同意请求则将好友存入数据库
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(2)
	defer dblink.Close()
}

func (o OperationRedisForUf) UserFriendRequestRefuse(friendUid int64) {
	panic("implement me")
}
