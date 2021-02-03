package db

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type DB interface {
	MySqlInit() *gorm.DB
	MySqlClose() error
	RedisInit(dbs int) *redis.Client
	RedisClose(client *redis.Client)
}

type SetData struct {
	Config *config.Config
}

func init()  {
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	sd := &SetData{}
	sd.Config = cond
}
