package model

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
)

type UserState interface {
	UserDelSignState(uid string) bool
	UserGetSignState(uid string) bool
	UserSearchSignUser(user interface{}) ([]int64,bool)
}

type OperationRedis struct {

}

func (o OperationRedis) UserDelSignState(uid string) bool {
	// 初始化redis连接
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer dblink.Close()
	_,err := dblink.Del(uid).Result()
	return err == redis.Nil
}

func (o OperationRedis) UserGetSignState(uid string) bool {
	// 初始化redis连接
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	defer dblink.Close()
	_,err := dblink.Get(uid).Result()
	// 返回成功/失败
	return err == redis.Nil
}

func (o OperationRedis) UserSearchSignUser(user interface{}) ([]int64, bool) {
	// 搜索登录用户
	// 初始化redis和mysql连接
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(1)
	var dd db.DB = &db.SetData{}
	dblink2 := dd.MySqlInit()
	defer func() {
		dblink2.Close()
		dblink.Close()
	}()
	switch reflect.TypeOf(user).Name() {
	case "string":
		// 根据用户名和昵称查询
		var newUsers = []int64{1}
		var us User = &Operations{}
		if userList,bl := us.SearchUser(user); bl {
			for i := 0;i < len(userList); i++ {
				if _,err := dblink.Get(strconv.FormatInt(userList[i], 10)).Result();err != redis.Nil {
					// 正确的结果在newUsers数组后面接着存储
					newUsers[len(newUsers) - 1] = userList[i]
				}
			}
			return newUsers,true
		}
	case "int64":
		// 根据id查询
		if _,err := dblink.Get(user.(string)).Result();err != redis.Nil {
			return []int64{user.(int64)}, false
		}
	}
	return nil, false
}
