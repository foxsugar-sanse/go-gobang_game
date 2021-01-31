package router

import (
	"github.com/gin-gonic/gin"
	"github.com/foxsuagr-sanse/go-gobang_game/app/controller"
)


type Router interface {
	Run(engine *gin.Engine)
}

type Route struct {

}

func (r *Route)Run(c *gin.Engine)  {
	// 实现关于user api部分的接口
	var user_op controller.RouterRequest = controller.UserRouter{}
	v1 := c.Group("/v1")
	{
		v1.POST("/user",user_op.LoginPost) // 提交一个新用户
		v1.GET("/user/:name",user_op.UserInfoGet) // 主要获取用户信息
		v1.PUT("/user",user_op.UserInfoUpdate) // 更新用户信息，需要鉴权
		v1.DELETE("/user",user_op.UserDelete) // 注销或者软删除用户
		v1.OPTIONS("/user",user_op.UserOtherOperations) // 用户接口的其他操作,比如search一个用户

		v1.GET("/rankles") // 根据name获取指定的排行榜，比如胜场榜，有num和page参数，num代表依次传回多少个用户，默认为10,page为分页

		v1.GET("/linkman") // 获取联系人列表，有num参数，代表一次返回多少个用户，默认为10
		v1.POST("/linkman/:name") // 有鉴权，将指定用户添加到联系人中，需要指定用户同意
		v1.DELETE("linkman/:name") // 鉴权，将指定用户重联系人中删除
		v1.PUT("linkman/:name") // 修改该联系人的信息，比如备注等等
		v1.OPTIONS("linkman/:name") // 联系人接口的其他功能，比如想联系人发出游戏邀请和信息等等

		// 历史记录接口不允许修改
		v1.GET("/history") // ?num=10&page=1
		v1.GET("/history/:history_id") // 鉴权接口，获取指定的一条历史记录的所有数据
		v1.DELETE("/history/:history_id") // 鉴权接口，删除指定的历史记录
	}
}