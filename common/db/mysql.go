package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
)

var Db *gorm.DB


func Init() {
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	// 初始化数据库连接
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",cond.ConfData.Mysql.Username,cond.ConfData.Mysql.Password,cond.ConfData.Mysql.Ipaddr,cond.ConfData.Mysql.Port,cond.ConfData.Mysql.Dbname)
	var err error
	Db, err = gorm.Open("mysql",dsn)
	if err != nil {
		panic(err)
	}
}

func Close() error {
	return Db.Close()
}