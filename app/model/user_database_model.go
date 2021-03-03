package model

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"reflect"
	"strconv"
	"github.com/ymzuiku/hit"
)

// 用户的数据库操作
type User interface {
	Login(user map[string]string) 									(bool, *errors.Errno)
	Sign(map[string]string) 										(int64,string,*errors.Errno)
	SetInfo(int64,map[string]interface{}) 							*errors.Errno
	GetUserInfo(uid int64) 											*Users
	DeleteUser(uid int64) 											*errors.Errno
	SearchUser(userSearch interface{}) 								([]int64,bool)
	AddUserFriend(uid int64,fid int64) 								bool
	SetUserFriendInfo(uid int64,fid int64,note string,group string) *errors.Errno
	DeleteUserFriend(uid int64,fid int64) 							bool
	QueryUserFriend(uid int64) 										([]*UserFriend,bool)
	FormGroupGetUserFriend(uid int64,group string) 					([]*UserFriend,bool)
	SetUserPortraitUrl(uid int64,url string)						*errors.Errno
	DeleteUserPortrait(uid int64)						            *errors.Errno
}

// 默认好友分组name
const DefaultUserGroupName = "我的好友"

// 数据库表映射 {user}
type Users struct {
	Id				int
	Uid 			int64
	UserName        string
	PassWord        string
	UserNickName    string
	UserAge         int
	UserSex         string
	UserBrief       string
	UserContact     string
	UserPortrait    string
	UserEmail       string
	UserPhone       int64
	EmailValidation int
	PhoneValidation int
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
	pwd.Write([]byte(user["UserPassWord"] + db.SALT + func() string {
		if len(User1) > 0 {
			s := strconv.FormatInt(db.GODETIME + int64(User1[len(User1) - 1].Id + 1),10)
			return s
		} else {
			s := strconv.FormatInt(db.GODETIME + int64(1),10)
			return s
		}

	}()))
	user["UserPassWord"] = hex.EncodeToString(pwd.Sum(nil))
	// 检测有没有重名用户
	var User []*Users
	dblink.Where("user_name = ?",user["UserName"]).First(&User)
	if len(User) > 0 {
		return false,errors.ErrUserNameExist
	} else {
		//err := dblink.Create(&Users{
		//	Uid: int64(2000 + len(User1) + 1),
		//	UserName: user["UserName"],
		//	PassWord: user["UserPassWord"],
		//}).Error
		//// 为用户创建一个默认分组{我的好友}
		//// {GroupRank : 0},0表示分组等级，0表示默认为系统创建用户时添加，用户添加的分组则按数据库规则默认为1
		//err2 := dblink.Create(&UserGroups{
		//	Uid:   int64(2000 + len(User1) + 1),
		//	Group: "我的好友",
		//	GroupRank: 0,
		//}).Error
		//return err == nil && err2 == nil,errors.OK

		// 使用事务确保两次操作成功
		err := dblink.Transaction(func(tx *gorm.DB) error {

			if err := tx.Create(&Users{
				Uid: int64(2000 + len(User1) + 1),
				UserName: user["UserName"],
				PassWord: user["UserPassWord"],
			}).Error; err != nil {
				return err
			}
			// 为用户创建一个默认分组{我的好友}
			// {GroupRank : 0},0表示分组等级，0表示默认为系统创建用户时添加，用户添加的分组则按数据库规则默认为1
			if err := tx.Create(&UserGroups{
				Uid:   int64(2000 + len(User1) + 1),
				Group: "我的好友",
				GroupRank: 0,
			}).Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return false,errors.ErrDatabase
		} else {
			return true,errors.OK
		}
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
		pwd.Write([]byte(user["UserPassWord"] + db.SALT + func() string {
			s := strconv.FormatInt(db.GODETIME + int64(User1[0].Id),10)
			return s
		}()))
		user["UserPassWord"] = hex.EncodeToString(pwd.Sum(nil))
	} else {
		return 0,"",errors.ErrUserNotFound
	}
	// 查询数据库
	var User []*Users
	dblink.Where("user_name = ? and pass_word = ?",user["UserName"],user["UserPassWord"]).First(&User)
	if len(User) > 0 {
		// 生成并返回token
		return User[0].Uid,User[0].UserName,errors.OK
	} else {
		return 0,"",errors.ErrUserNotFound
	}
}

func (op *Operations)SetInfo(uid int64,user map[string]interface{}) *errors.Errno {
	// TODO:头像等设置还未做
	// 初始化数据库连接
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()

	// 根据uid查询具体要设置的用户信息
	var User []*Users
	dblink.Where("uid = ?",uid).Find(&User)
	if len(User) > 0 {
		dblink.Model(&Users{}).Where("uid = ?",uid).Updates(map[string]interface{}{
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
		var User3 []*Users
		dblink.Where("user_name = ?",userSearch).First(&User2)
		dblink.Where("user_nick_name = ?",userSearch).Find(&User3)
		if len(User2) > 0 {
			var userSlice = make([]int64,1)
			userSlice[0] = User2[0].Uid
			// 去重
			for i := 0; i < len(User3); i++ {
				if User2[0].Uid == User3[i].Uid{
					// 重复不做任何操作
				} else {
					// 不重复则对切片{userSlice}扩容
					userSlice = append(userSlice,User3[i].Uid)
				}
			}
			return userSlice,true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (op *Operations) AddUserFriend(uid int64, fid int64) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserF []*UserFriend
	// 查询有没有重复的好友关系
	dblink.Where("main_uid = ? and friend_uid = ?",uid,fid).Find(&UserF)
	if len(UserF) > 0 {
		return false
	} else {
		dblink.Create(&UserFriend{
			MainUid:   uid,
			FriendUid: fid,
			UserGroup: DefaultUserGroupName,
		})
		return true
	}
}

func (op *Operations) SetUserFriendInfo(uid int64, fid int64, note string, group string) *errors.Errno {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserF []*UserFriend
	// 查询是否为好友
	dblink.Where("main_uid = ? and friend_uid = ?",uid,fid).First(&UserF)
	if len(UserF) > 0 {
		var UserGroup []*UserGroups
		// 检查用户设置的分组在数据库中有没有存在
		dblink.Where("uid = ? and group = ?",uid,group).First(&UserGroup)
		if len(UserGroup) == 0 {
			return errors.ErrGroupNotFound
		}
		dblink.Model(&UserF).Where("main_uid = ? and friend_uid = ?",uid,fid).Updates(&UserFriend{
			FriendNote: note,
			UserGroup: group,
		})
		return errors.OK
	} else {
		return errors.ErrBadRequest
	}
}

func (op *Operations) DeleteUserFriend(uid int64, fid int64) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserF []*UserFriend
	// 查询要删除的用户
	dblink.Where("main_uid = ? and friend_uid = ?",uid,fid).First(&UserF)
	if len(UserF) > 0 {
		dblink.Where("main_uid = ? and friend_uid = ?",uid,fid).Delete(&UserF)
		return true
	} else {
		return false
	}
}

func (op *Operations) QueryUserFriend(uid int64) ([]*UserFriend,bool) {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserF []*UserFriend
	// 查询好友关系
	dblink.Where("main_uid = ?",uid).Find(&UserF)
	if len(UserF) > 0 {
		return UserF, true
	} else {
		return nil, false
	}
}

func (op Operations) FormGroupGetUserFriend(uid int64, group string) ([]*UserFriend,bool) {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserGroup []*UserGroups
	// 检查用户设置的分组在数据库中有没有存在
	dblink.Where("uid = ? and group = ?", uid, group).First(&UserGroup)
	if len(UserGroup) > 0 {
		var UserFriends []*UserFriend
		dblink.Where("main_uid = ? and user_group = ?", uid, group).Find(&UserFriends)
		return UserFriends, len(UserFriends) > 0
	} else {
		return nil, false
	}
}

func (op *Operations) SetUserPortraitUrl(uid int64, url string) *errors.Errno {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var user []*Users
	dblink.Where("uid = ?",uid).First(&user)
	if len(user) > 0 {
		// 如果已有图片则删除覆盖
		if user[0].UserPortrait != "." && user[0].UserPortrait != "" {
			var con config.ConFig = &config.Config{}
			conf := con.InitConfig()
			err := os.Remove(conf.ConfData.Model.Localurl + "/" + user[0].UserPortrait)
			if err != nil {
				log.Println(err)
			}
		}
		dblink.Model(&user).Where("uid = ?",uid).Updates(&Users{
			UserPortrait: url,
		})
		return errors.UploadOK
	} else {
		return errors.ErrUserNotFound
	}
}

func (op *Operations) DeleteUserPortrait(uid int64) *errors.Errno {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var user []*Users
	dblink.Where("uid = ?",uid).First(&user)
	if len(user) > 0 {
		// 被删除的头像路径一律改为"."
		if err := dblink.Model(&user).Where("uid = ?",uid).Updates(&Users{
			UserPortrait: ".",
		}).Error ; err == nil {
			var con config.ConFig = &config.Config{}
			conf := con.InitConfig()
			err2 := os.Remove(conf.ConfData.Model.Localurl + "/" + user[0].UserPortrait)
			return hit.If(err2 == nil,errors.OK,errors.ErrUserUploadUrlNot).(*errors.Errno)
		} else {
			return errors.ErrDatabase
		}
	} else {
		return errors.ErrUserNotFound
	}
}

