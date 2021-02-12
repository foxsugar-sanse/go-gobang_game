package model

import (
	"github.com/foxsuagr-sanse/go-gobang_game/common/db"
)

type UserGroupDataBase interface {
	// crud接口
	AddUserGroup(uid int64,group string) bool
	DeleteUserGroup(uid int64,group string) bool
	SetUserGroup(uid int64,group string,setGroup string) bool
	GetUserGroup(uid int64) []string

}

// 数据库表映射
type UserGroup struct {
	Id 				int
	Uid 			int64
	Group 			string
	GroupRank 		int
}

type OperationsUserGroup struct {

}

func (o OperationsUserGroup) AddUserGroup(uid int64, group string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	// 检测有没有相同
	if func() bool {
		var UserGroups []*UserGroup
		dblink.Where("uid = ? And group = ?",uid,group).First(&UserGroups)
		return len(UserGroups) == 0
	}() {
		// 没有相同则添加
		err := dblink.Model(&UserGroup{}).Create(&UserGroup{
			Uid:   uid,
			Group: group,
		}).Error
		return err == nil
	} else {
		return false
	}
}

func (o OperationsUserGroup) DeleteUserGroup(uid int64, group string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	if func() bool {
		var UserGroups []*UserGroup
		dblink.Where("uid = ? And group = ?",uid,group).First(&UserGroups)
		return len(UserGroups) == 1
	}() {
		err2 := dblink.Model(&UserGroup{}).Where("uid = ? And group = ?",uid,group).Delete(&UserGroup{}).Error
		return err2 == nil
	} else {
		return false
	}
}

func (o OperationsUserGroup) SetUserGroup(uid int64, group string, setGroup string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	if func() bool {
		var UserGroups []*UserGroup
		dblink.Where("uid = ? And group = ?",uid,group).First(&UserGroups)
		return len(UserGroups) == 1
	}() {
		err2 := dblink.Model(&UserGroup{}).Where("uid = ? And group = ?",uid,group).Updates(&UserGroup{
			Group: setGroup,
		})
		return err2 == nil
	} else {
		return false
	}
}

func (o OperationsUserGroup) GetUserGroup(uid int64) []string {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserGroups []*UserGroup
	dblink.Where("uid = ?",uid).Find(&UserGroups)
	strSlice := make([]string,0)
	for i := 0 ; i < len(UserGroups); i++ {
		strSlice = append(strSlice,UserGroups[i].Group)
	}
	return strSlice
}


