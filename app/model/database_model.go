package model

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	//"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
)

// 用户的数据库操作
type User interface {
	Login(map[string]string) (bool, int)
	Sign(map[string]string) (string,int)
	SetInfo(map[string]interface{}) (string, int)
	GetUserInfo(uid int) (string, int)
}

// 数据库表映射
type Users struct {
	Uid 			int
	UserName 		string
	UserNickName 	string
	UserAge			string
	UserSex			string
	UserBrief		string
	UserContact		string
}

type Operations struct {

}


func (op *Operations) Login(user map[string]string) (bool,int) {
	db.Init()
	db.Db.Create(&Users{})
	return false,0
}

func (op *Operations) Sign(user map[string]string) (string,int){
	return "",0
}

func (op *Operations)SetInfo(user map[string]interface{}) (string, int) {
	return "",0
}

func (op *Operations)GetUserInfo(uid int) (string, int) {

	return "", 200
}
