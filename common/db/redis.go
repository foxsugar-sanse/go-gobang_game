package db

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/go-redis/redis"
)



func (sd * SetData) RedisInit(dbs int) *redis.Client {
	// 连接redis数据库
	var con config.ConFig = &config.Config{}
	conf := con.InitConfig()
	Client := redis.NewClient(&redis.Options{
		Addr: conf.ConfData.Redis.Ipaddr + conf.ConfData.Redis.Port,
		Password: conf.ConfData.Redis.Password,
		DB: dbs,
	})
	_,err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return Client
}

func (sd * SetData) RedisClose(client *redis.Client)  {
	_ = client.Close()
}