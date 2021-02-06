package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/foxsuagr-sanse/go-gobang_game/app/controller"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	//"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"reflect"
	"strconv"
)

// 用户的数据库操作
type User interface {
	Login(user map[string]string) (bool, *errors.Errno)
	Sign(map[string]string) (int64,string,*errors.Errno)
	SetInfo(int64,map[string]interface{}) *errors.Errno
	GetUserInfo(uid int64) *Users
	DeleteUser(uid int64) *errors.Errno
	SearchUser(userSearch interface{}) ([]int64,bool)
}

// 数据库表映射 {user}
type Users struct {
	Id				int
	Uid 			int64
	UserName 		string
	PassWord		string
	UserNickName 	string
	UserAge			string
	UserSex			string
	UserBrief		string
	UserContact		string
}

// 数据库表映射 {user_friend}
type UserFriend struct {
	Id 				int
	MainUid			int64
	FriendUid		int64
	FriendNote		string
	UserGroup 		string
}

type Operations struct {

}


func (op *Operations) Login(user map[string]string) (bool,*errors.Errno) {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	// 密码加密策略,通过规定的时间戳，数据表里有多少用户字段就读取然后len加到原来的时间戳和一个固定的key组成盐值
	var User1 []*Users
	dblink.Find(&User1)
	pwd := md5.New()
	pwd.Write([]byte(controller.SALT + func() string {
		s := strconv.FormatInt(controller.GODETIME + int64(len(User1)),10)
		return s
	}()))
	user["UserPassWord"] = hex.EncodeToString(pwd.Sum(nil))
	// 检测有没有重名用户
	var User []*Users
	dblink.Where("user_name = ?",user["UserName"]).First(&User)
	if len(User) > 0 {
		return false,errors.ErrUserNameExist
	} else {
		dblink.Create(&Users{
			Uid: 60014,
			UserName: user["UserName"],
			PassWord: user["UserPassWord"],
		})
		return true,errors.OK
	}
}

func (op *Operations) Sign(user map[string]string) (int64,string,*errors.Errno){
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	// 登录密码加密策略，根据用户名到数据库中获取指定用户的id，(id + 时间戳 + 固定字符串值)组成盐值
	var User1 []*Users
	dblink.Where("user_name = ?",user["UserName"]).First(&User1)
	if len(User1) > 0 {
		pwd := md5.New()
		pwd.Write([]byte(controller.SALT + func() string {
			s := strconv.FormatInt(controller.GODETIME + int64(User1[0].Id),10)
			return s
		}()))
		user["UserPassWord"] = hex.EncodeToString(pwd.Sum(nil))
	} else {
		return 0,"",errors.ErrUserNotFound
	}
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

func (op *Operations) SearchUser(userSearch interface{}) ([]int64,bool) {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	switch reflect.TypeOf(userSearch).Name() {
	case "int64":
		// 根据uid进行查询
		var User []*Users
		dblink.Where("uid = ?",userSearch).First(&User)
		if len(User) > 0 {
			// 查到匹配的用户
			return []int64{User[0].Uid},true
		} else {
			return nil, false
		}
	case "string":
		// 根据用户名和昵称进行查询
		var User2 []*Users
		dblink.Where("user_name = ?",userSearch).First(&User2)
		dblink.Where("user_nick_name = ?",userSearch).Find(&User2)
		if len(User2) > 0 {
			var userSlice = []int64{1}
			for i := 0;i < len(User2);i++ {
				userSlice[i] = User2[i].Uid
			}
			return userSlice,true
		} else {
			return nil, false
		}
	}
	return nil, false
}