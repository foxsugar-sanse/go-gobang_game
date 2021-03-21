package model

import (
	"github.com/bwmarrin/snowflake"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"log"
	"strconv"
)

/*
	5数据库存储待绑定的客户端id
*/

const RedisDBC int = 5

type ClientState interface {
	New() (string, error)
	Bind(cid string, uid string) error
	Delete(cid string) error
}

type Clients struct{}

func (c Clients) New() (string, error) {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(RedisDBC)
	defer dblink.Close()

	// 生成id并返回
	// 雪花算法
	n, err := snowflake.NewNode(0)
	if err != nil {
		log.Println(err)
		return "", err
	}
	id := n.Generate()
	cid := strconv.FormatInt(int64(id), 10)
	// 存入redis
	dblink.Set(cid, "", 0)
	return cid, nil
}

func (c Clients) Bind(cid string, uid string) error {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(RedisDBC)
	defer dblink.Close()

	// 将客户端id与用户uid绑定
	return dblink.Set(cid, uid, 0).Err()
}

func (c Clients) Delete(cid string) error {
	var d db.DB = &db.SetData{}
	dblink := d.RedisInit(RedisDBC)
	defer dblink.Close()

	// 删除客户端id与其绑定的值
	return dblink.Del(cid).Err()
}
