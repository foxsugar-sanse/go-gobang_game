package router

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/controller"
	"github.com/foxsuagr-sanse/go-gobang_game/router/middleware"
	"github.com/gin-gonic/gin"
)


type Router interface {
	Run(engine *gin.Engine)
}

type Route struct {

}

func (r *Route)Run(c *gin.Engine)  {
	// 实现关于user api部分的接口
	var userOp controller.RouterRequest = &controller.UserRouter{}
	v1 := c.Group("/v1")
	v1.Use(middleware.JwtMiddlewareOAuth())
	{
		v1.GET("/user", userOp.UserInfoGet)           		        // 主要获取用户信息,?uid=
		v1.GET("/user/sign", userOp.UserGetSignState)       		// 查询用户是否登录,?uid=
		v1.PUT("/user", userOp.UserInfoUpdate)              		// 更新用户信息，需要鉴权
		v1.DELETE("/user/sign", userOp.UserDeleteSignState) 		// 注销登录用户
		v1.DELETE("/user", userOp.UserDelete)               		// 软删除用户
		v1.OPTIONS("/user", userOp.UserOtherOperations)     		// 用户接口的其他操作,比如search一个用户

		v1.GET("/rankles") 										// 根据name获取指定的排行榜，比如胜场榜，有num和page参数，num代表依次传回多少个用户，默认为10,page为分页

		v1.GET("/friend_request",userOp.GetUserFriendRequest)		// 获取好友申请
		v1.POST("/friend_request",userOp.ConsentUserFriendRequest) // 同意好友申请
		v1.DELETE("/friend_request",userOp.RefuseUserFriendRequest)// 拒绝好友申请
		v1.PUT("/friend_request",userOp.CreateUserFriendRequest)   // 创建好友申请
		v1.GET("/linkman",userOp.GetUserForFriend) 				// 获取联系人列表，有num参数，代表一次返回多少个用户，默认为10
		//v1.POST("/linkman/:name",userOp.AddUserForFriend) 			// 有鉴权，将指定用户添加到联系人中，需要指定用户同意
		v1.DELETE("linkman/:name",userOp.DeleteUserForFriend) 		// 鉴权，将指定用户重联系人中删除
		v1.PUT("linkman/:name",userOp.ModifyFriendInfo) 			// 修改该联系人的信息，比如备注等等
		v1.OPTIONS("linkman/:name",userOp.OtherUserFriendInterface)// 联系人接口的其他功能，比如想联系人发出游戏邀请和信息等等

		// 历史记录接口不允许修改
		v1.GET("/history") 										// ?num=10&page=1
		v1.GET("/history/:history_id") 							// 鉴权接口，获取指定的一条历史记录的所有数据
		v1.DELETE("/history/:history_id") 							// 鉴权接口，删除指定的历史记录
	}
	// 公共接口不需要鉴权
	v1Pub := c.Group("/v2")
	v1Pub.Use(middleware.TimeMatch())
	{
		//v1Pub.POST("/user",userOp.LoginPost) 									// 提交一个新用户
		//v1Pub.GET("/user/login")              								// 获取注册所需的一些参数
		v1Pub.POST("/user/login", userOp.LoginPost) 				// 注册用户
		//v1Pub.GET("/user/sign")               								// 获取登录所需的一些参数
		v1Pub.POST("/user/sign", userOp.SignPost) 					// 创建用户登录
	}

}