package model

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	//"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
)

// 用户的数据库操作
type User interface {
	Login(map[string]string) (bool, *errors.Errno)
	Sign(map[string]string) (int64,string,*errors.Errno)
	SetInfo(int64,map[string]interface{}) *errors.Errno
	GetUserInfo(uid int64) *Users
	DeleteUser(uid int64) *errors.Errno
}

// 数据库表映射
type Users struct {
	Uid 			int64
	UserName 		string
	PassWord		string
	UserNickName 	string
	UserAge			string
	UserSex			string
	UserBrief		string
	UserContact		string
}


type Operations struct {

}


func (op *Operations) Login(user map[string]string) (bool,*errors.Errno) {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()

	// 检测有没有重名用户
	var User []*Users
	dblink.Where("user_name = ?",user["UserName"]).First(&User)
	if len(User) > 0 {
		return false,errors.ErrUserNameExist
	} else {
		dblink.Create(&Users{
			Uid: 60014,
			UserName: user["UserName"],
			PassWord: user["UserPassWod"],
		})
		return true,errors.OK
	}
}

func (op *Operations) Sign(user map[string]string) (int64,string,*errors.Errno){
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	// 查询数据库
	var User []*Users
	dblink.Where("user_name = ? and user_password = ?",user["UserName"],user["UserPassWord"]).First(&User)
	if len(User) > 1 {
		// 生成并返回token
		return User[0].Uid,User[0].UserName,errors.OK
	} else {
		return 0,"",errors.ErrUserNotFound
	}
}

func (op *Operations)SetInfo(uid int64,user map[string]interface{}) *errors.Errno {
	// 初始化数据库连接
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()

	// 根据uid查询具体要设置的用户信息
	var User []*Users
	dblink.Where("uid = ?",uid).First(&User)
	if len(User) > 0 {
		dblink.Model(&User).Updates(map[string]interface{}{
			"user_nick_name":user["UserNickName"],
			"user_brief":user["UserBrief"],
			"user_age":user["UserAge"],
			"user_sex":user["UserSex"],
			"user_contact":user["UserContact"],
		})
		return errors.OK
	} else {
		return errors.ErrUserNotFound
	}
}

func (op *Operations)GetUserInfo(uid int64) *Users {
	// 查询用户信息
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()

	var User []*Users
	dblink.Where("uid = ?",uid).First(&User)
	return User[0]
}

func (op *Operations) DeleteUser(uid int64) *errors.Errno {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()

	// 根据uid查询要删除的用户
	var User []*Users
	dblink.Where("uid = ?",uid).First(&User)
	if len(User) > 0 {
		dblink.Where("uid = ?",uid).Delete(&User)
		return errors.OK
	} else {
		return errors.ErrUserNotFound
	}
}
