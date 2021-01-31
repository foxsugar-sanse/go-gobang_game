package db

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
)

// 用户的数据库操作
type User interface {
	Login(map[string]string) (bool, int)
	Sign(map[string]string) (string,int)
	SetInfo(map[string]interface{}) (string, int)
	GetUserInfo(uid int) (string, int)
}

type Operations struct {}

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

func (op *Operations) Login(user map[string]string) (bool,int) {

}

func (op *Operations) Sign(user map[string]string) (string,int){

}

func (op *Operations)SetInfo(user map[string]interface{}) (string, int) {

}

func (op *Operations)GetUserInfo(uid int) (string, int) {

	return "", 200
}