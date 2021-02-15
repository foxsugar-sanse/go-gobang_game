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
type UserGroups struct {
	Id 				int
	Uid 			int64
	Group 			string `gorm:"column:user_group"`
	GroupRank 		int
}

type OperationsUserGroup struct {

}

func (o OperationsUserGroup) AddUserGroup(uid int64, group string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	// 检测有没有相同
	if bl := func() bool {
		var UserGroup []*UserGroups
		dblink.Where("uid = ? and user_group = ?",uid,group).First(&UserGroup)
		return len(UserGroup) == 0
	}() ; bl {
		// 没有相同则添加
		err := dblink.Model(&UserGroups{}).Create(&UserGroups{
			Uid:   uid,
			Group: group,
			GroupRank: 1,
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
	if bl := func() bool {
		var UserGroup []*UserGroups
		dblink.Where("uid = ? and user_group = ? and group_rank = ? ",uid,group,1).Find(&UserGroup)
		return len(UserGroup) == 1
	}() ; bl{
		err2 := dblink.Model(&UserGroups{}).Where("uid = ? and user_group = ?",uid,group).Delete(&UserGroups{}).Error
		return err2 == nil
	} else {
		return false
	}
}

func (o OperationsUserGroup) SetUserGroup(uid int64, group string, setGroup string) bool {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserGroup []*UserGroups
	if bl := func() bool {
		dblink.Where("uid = ? and user_group = ? and group_rank = ?",uid,group,1).Find(&UserGroup)
		return len(UserGroup) != 0
	}(); bl {
		err2 := dblink.Model(&UserGroups{}).Where("uid = ? and user_group = ?",uid,group).Updates(&UserGroups{
			Group: setGroup,
		}).Error
		return err2 == nil
	} else {
		return false
	}
}

func (o OperationsUserGroup) GetUserGroup(uid int64) []string {
	var d db.DB = &db.SetData{}
	dblink := d.MySqlInit()
	defer dblink.Close()
	var UserGroups []*UserGroups
	dblink.Where("uid = ?",uid).Find(&UserGroups)
	strSlice := make([]string,0)
	for i := 0 ; i < len(UserGroups); i++ {
		strSlice = append(strSlice,UserGroups[i].Group)
	}
	return strSlice
}


