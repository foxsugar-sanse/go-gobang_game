package model

import (
	"encoding/json"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"strconv"
)


type UserForFriend interface {
	UserFriendRequestCreate(mainUid int64,friendUid int64,note string) 	bool
	UserFriendRequestGet(friendUid int64) 								(map[string]map[string]string,bool)
	UserFriendRequestConsent(friendUid int64,mainUid int64) 			bool
	UserFriendRequestRefuse(friendUid int64,mainUid int64)				bool
}

type OperationRedisForUf struct {

}

func (o OperationRedisForUf) UserFriendRequestCreate(mainUid int64, friendUid int64, note string) bool{
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer dblink.Close()
	MainUid := strconv.FormatInt(mainUid,10)
	// key格式"friendUid"
	// key-value 存储多个好友申请，string 为未序列化的json字符串
	// 获取Redis中有没有已存储记录，防止覆盖
	if ResMapStr,err := dblink.Get(strconv.FormatInt(friendUid,10)).Result();err == nil {
		ResMap := make(map[string]map[string]string)
		_ = json.Unmarshal([]byte(ResMapStr),&ResMap)
		for i := 0; i < len(ResMapStr);i++{
			if ResMap[MainUid]["state"] == "" {
				// 为空则添加
				ResMap[MainUid] = make(map[string]string)
				//ResMap[MainUid]["main_uid"] = strconv.FormatInt(mainUid,10)
				//ResMap[MainUid]["friend_uid"] = strconv.FormatInt(friendUid,10)
				ResMap[MainUid]["note"] = note
				ResMap[MainUid]["state"] = "0"
				mapJSONRes,_ := json.Marshal(ResMap)
				err2 := dblink.Set(strconv.FormatInt(friendUid, 10),mapJSONRes,0).Err()
				return err2 == nil
			}
		}
	} else {
		// 初始化map
		ResMap2 := make(map[string]map[string]interface{})
		ResMap2[MainUid] = make(map[string]interface{})
		ResMap2[MainUid]["note"] = note
		ResMap2[MainUid]["state"] = "0"
		mapJSONRes,_ := json.Marshal(ResMap2)
		err2 := dblink.Set(strconv.FormatInt(friendUid, 10),mapJSONRes,0).Err()
		return err2 == nil
	}
	//err := dblink.HMSet(strconv.FormatInt(friendUid, 10), map[string]interface{}{
	//	"main_uid": mainUid,
	//	"friend_uid": friendUid,
	//	"note": note,
	//	"state":0,
	//}).Err()
	return false
}

func (o OperationRedisForUf) UserFriendRequestGet(friendUid int64) (map[string]map[string]string,bool) {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer dblink.Close()
	if requestMapStr, err := dblink.Get(strconv.FormatInt(friendUid, 10)).Result();err == nil {
		// 返回一个json序列化之后的map
		ResMap := map[string]map[string]string{}
		_ = json.Unmarshal([]byte(requestMapStr),&ResMap)
		return ResMap,true
	} else {
		return nil,false
	}
}

func (o OperationRedisForUf) UserFriendRequestConsent(friendUid int64,mainUid int64) bool {
	// 同意请求则将好友存入数据库
	// TODO 多个好友申请未解决！
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer func() {
		dblink.Close()
	}()
	MainUid := strconv.FormatInt(mainUid,10)
	if ResMapStr,err := dblink.Get(strconv.FormatInt(friendUid, 10)).Result();err == nil {
		// 好友申请存在则删除
		// 匹配好友申请并删除
		ResMap := map[string]map[string]string{}
		_ = json.Unmarshal([]byte(ResMapStr),&ResMap)
		for i := 0; i < len(ResMap);i++ {
			if ResMap[MainUid]["state"] == "0" {
				// 匹配成功
				// 长度为1直接删除
				if len(ResMap) == 1 {
					_ = dblink.Del(strconv.FormatInt(friendUid, 10)).Err()
				} else {
					// 不为1则使用delete函数
					delete(ResMap,MainUid)
					mapJSONRes,_ := json.Marshal(ResMap)
					dblink.Set(strconv.FormatInt(friendUid, 10),mapJSONRes,0)
				}
				var u User = &Operations{}
				return u.AddUserFriend(mainUid,friendUid)
			} else {

			}
		}

	} else {
		return false
	}

	return false
}

func (o OperationRedisForUf) UserFriendRequestRefuse(friendUid int64, mainUid int64) bool {
	// 拒绝申请直接从Redis数据库中删除申请
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer func() {
		dblink.Close()
	}()
	MainUid := strconv.FormatInt(mainUid,10)
	// 获取数据
	if ResMapStr,err := dblink.Get(strconv.FormatInt(friendUid, 10)).Result();err == nil {
		ResMap := map[string]map[string]string{}
		_ = json.Unmarshal([]byte(ResMapStr),&ResMap)
		for i := 0; i < len(ResMap); i++ {
			if ResMap[MainUid]["state"] == "0" {
				// 删除并重新设置
				delete(ResMap,MainUid)
				mapJSONRes,_ := json.Marshal(ResMap)
				dblink.Set(strconv.FormatInt(friendUid, 10),mapJSONRes,0)
				return true
			} else {
				//
			}
		}
	}
	return false
}
