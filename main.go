package main

import (
	"fmt"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/router"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 读取配置文件
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	r := gin.Default()
	// 初始化路由
	var x router.Router = &router.Route{}
	x.Run(r)
	// 初始化数据库连接
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",cond.ConfData.Mysql.Username,cond.ConfData.Mysql.Password,cond.ConfData.Mysql.Ipaddr,cond.ConfData.Mysql.Port,cond.ConfData.Mysql.Dbname)
	db, err := gorm.Open("mysql",dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_ = r.Run(cond.ConfData.Run.Ipaddr + cond.ConfData.Run.Port)
}
