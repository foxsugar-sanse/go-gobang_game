package model

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/go-redis/redis"
	"strconv"
)

type InviteRedis interface{
	CreateInvite(key string,value string) bool
	GetInvite(key string) ([]string,bool)
	DeleteInvite(key string,mid int64) bool
}

type Invite struct {}

func (i Invite) CreateInvite(key string, value string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(3)
	defer dblink.Close()
	// 格式{2001+0}{uid+num}
	// 获取有无相同key
	for i := 0; i < 10; i++ {
		if _,err := dblink.Get(key + "+" + strconv.Itoa(i)).Result() ; err == redis.Nil {
			return dblink.Set(key + "+" + strconv.Itoa(i),value,3000).Err() == nil
		}
	}
	return dblink.Set(key + "+" + "0",value,3000).Err() == nil
}

func (i Invite) GetInvite(key string) ([]string, bool) {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(3)
	defer dblink.Close()
	slice := make([]string, 0)
	for i := 0; i < 10; i++ {
		if res,err := dblink.Get(key + "+" + strconv.Itoa(i)).Result() ; err != redis.Nil {
			slice = append(slice,res)
		}
	}
	return slice,len(slice) > 0
}

func (i Invite) DeleteInvite(key string,mid int64) bool {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(3)
	defer dblink.Close()
	return dblink.Del(key).Err() == nil
}
