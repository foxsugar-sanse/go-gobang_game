package model

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/go-redis/redis"
	"strconv"
)


type UserForFriend interface {
	UserFriendRequestCreate(mainUid int64,friendUid int64,note string) 	bool
	UserFriendRequestGet(friendUid int64) 								(map[string]string,bool)
	UserFriendRequestConsent(friendUid int64) 							bool
	UserFriendRequestRefuse(friendUid int64)							bool
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

func (o OperationRedisForUf) UserFriendRequestConsent(friendUid int64) bool {
	// 同意请求则将好友存入数据库
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(2)
	defer func() {
		dblink.Close()
	}()
	if userRq,err := dblink.HGetAll(strconv.FormatInt(friendUid, 10)).Result();err != redis.Nil {
		// 好友申请存在则删除
		dblink.HDel(strconv.FormatInt(friendUid, 10))
		var u User = &Operations{}
		mainUid,_ := strconv.ParseInt(userRq["main_uid"],10,64)
		friendUid,_ := strconv.ParseInt(userRq["friend_uid"],10,64)
		return u.AddUserFriend(mainUid,friendUid)
	} else {
		return false
	}

}

func (o OperationRedisForUf) UserFriendRequestRefuse(friendUid int64) bool {
	// 拒绝申请直接从Redis数据库中删除申请
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(2)
	defer func() {
		dblink.Close()
	}()
	// 删除成功与否
	_,err := dblink.HDel(strconv.FormatInt(friendUid, 10)).Result()
	return err != nil
}
