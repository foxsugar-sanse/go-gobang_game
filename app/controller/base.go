package controller

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"strings"
)

const SALT string = "0x23&&%%GWGwyn12"

const GODETIME int64 = 1612428719

type RouterRequest interface {
	LoginPost				(c * gin.Context)
	UserInfoGet				(c * gin.Context)
	UserInfoUpdate			(c * gin.Context)
	UserDelete				(c * gin.Context)
	UserOtherOperations		(c * gin.Context)
	SignPost				(c * gin.Context)
	UserDeleteSignState		(c * gin.Context)
	UserGetSignState		(c * gin.Context)

	GetUserForFriend		(c * gin.Context)
	AddUserForFriend		(c * gin.Context)
	DeleteUserForFriend		(c * gin.Context)
	ModifyFriendInfo		(c * gin.Context)
	OtherUserFriendInterface(c * gin.Context)
	RefuseUserFriendRequest	(c * gin.Context)
	ConsentUserFriendRequest(c * gin.Context)
	GetUserFriendRequest	(c * gin.Context)
	CreateUserFriendRequest	(c * gin.Context)
}

type UserBindJsonOtherOpera struct {
	Opera			string `json:"opera"`
	UserBindJsonOtherOperaDataEntity `json:"data"`
}

type UserBindJsonOtherOperaDataEntity struct {
	Uid 			int64  `json:"uid"`
	UserName		string `json:"user_name"`
	UserNickName	string `json:"user_nick_name"`
}

type UserJsonBindLogin struct {
	TimeSlices 		int64  `json:"t_s"`
	UserName 		string `json:"user_name"`
	UserPassword 	string `json:"user_password"`
}

type UserJsonBindSign struct {
	TimeSlices 		int64  `json:"t_s"`
	UserName		string `json:"user_name"`
	UserPassWord	string `json:"user_password"`
}

type UserJsonBindUpdate struct {
	UserNickName	string `json:"user_nick_name"`
	UserBrief		string `json:"user_brief"`
	UserSex			string `json:"user_sex"`
	UserAge			string `json:"user_age"`
	UserContact		string `json:"user_contact"`
}

type UserRouter struct {
	UserJBL *UserJsonBindLogin
	UserJBS *UserJsonBindSign
}

func (u *UserRouter) LoginPost(c *gin.Context) {
	json := UserJsonBindLogin{}
	err := c.ShouldBindBodyWith(&json,binding.JSON)
	if err != nil {
		panic(err)
	}

	var md model.User = &model.Operations{}
	_, Opcode := md.Login(map[string]string{
		"UserName":json.UserName,
		"UserPassWord":json.UserPassword,
	})
	c.JSON(Opcode.HttpCode,gin.H{
		"code":    Opcode.Code,
		"message": Opcode.Message,
	})
}

func (u *UserRouter) UserInfoGet(c *gin.Context) {
	// 根据参数获取用户信息
	var md model.User = &model.Operations{}
	idArgs, _ := c.Params.Get("uid")
	uid,_ := strconv.ParseInt(idArgs,10,64)
	users := md.GetUserInfo(uid)
	c.JSON(errors.OK.HttpCode,gin.H{
		"code":errors.OK.Code,
		"message":errors.OK.Message,
		"data":map[string]interface{}{
			"uid":            users.Uid,
			"user_name":      users.UserName,
			"user_nick_name": users.UserNickName,
			"user_brief":     users.UserBrief,
			"user_age":       users.UserAge,
			"user_sex":       users.UserSex,
			"user_contact":   users.UserContact,
		},
	})
}

func (u *UserRouter) UserInfoUpdate(c *gin.Context) {
	// 绑定提交json
	json := UserJsonBindUpdate{}
	_ = c.BindJSON(&UserJsonBindUpdate{})
	// 鉴权通过，获取token信息
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead," ",2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	MyClaims,Op := jwt.MatchToken(tokenHeadInfo[1])
	if Op {
		// token解密成功
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(MyClaims.Uid,10,64)
		code := md.SetInfo(uid,map[string]interface{}{
			"UserNickName":json.UserNickName,
			"UserContact":json.UserContact,
			"UserSex":json.UserSex,
			"UserAge":json.UserAge,
			"UserBrief":json.UserBrief,
		})
		if code.Code == 0 {
			// 设置成功响应
			c.JSON(code.HttpCode,gin.H{
				"code":code.Code,
				"message":code.Message,
			})
		} else {
			// 设置失败响应
			c.JSON(code.HttpCode,gin.H{
				"code":code.Code,
				"message":code.Message,
			})
		}
	}
}

func (u *UserRouter) UserDelete(c *gin.Context) {
	// 鉴权通过，获取用户token，得知要删除的对应用户
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead," ",2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	MyClaims,Op := jwt.MatchToken(tokenHeadInfo[1])
	if Op {
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(MyClaims.Uid,10,64)
		code := md.DeleteUser(uid)
		c.JSON(code.HttpCode,gin.H{
			"code":code.Code,
			"message":code.Message,
		})
	}
}

func (u *UserRouter) UserOtherOperations(c *gin.Context) {
	json := &UserBindJsonOtherOpera{}
	_ = c.BindJSON(&UserBindJsonOtherOpera{})

	switch json.Opera {
	case "search_u":
		// 搜索一个用户
		if json.UserBindJsonOtherOperaDataEntity.Uid != 0 {
			var d model.User = &model.Operations{}
			userList,bl := d.SearchUser(json.UserBindJsonOtherOperaDataEntity.Uid)
			if bl {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
					"data":userList[0],
				})
			} else {
				c.JSON(errors.ErrUserNotFound.HttpCode,gin.H{
					"code":errors.ErrUserNotFound.Code,
					"message":errors.ErrUserNotFound.Message,
				})
			}
		} else if json.UserBindJsonOtherOperaDataEntity.UserName != "" {
			// 昵称或者用户名搜索
			var d model.User = &model.Operations{}
			userList,bl := d.SearchUser(json.UserBindJsonOtherOperaDataEntity.UserName)
			userMap := make(map[int]int64)
			if bl {
				for i := 1;i <= len(userList);i++ {
					userMap[1] = userList[i-1]
				}
				// 返回用户数据
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
					"data":userMap,
				})
			} else {
				c.JSON(errors.ErrUserNotFound.HttpCode,gin.H{
					"code":errors.ErrUserNotFound.Code,
					"message":errors.ErrUserNotFound.Message,
				})
			}
		}
	case "search_p":
		// 搜索一个在线用户
		if json.UserBindJsonOtherOperaDataEntity.Uid != 0 {
			var d model.UserState = &model.OperationRedis{}
			if userList,bl := d.UserSearchSignUser(json.UserBindJsonOtherOperaDataEntity.Uid); bl {
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
					"data":userList[0],
				})
			}
		} else if json.UserBindJsonOtherOperaDataEntity.UserName != "" {
			var d model.UserState = &model.OperationRedis{}
			if userList,bl := d.UserSearchSignUser(json.UserBindJsonOtherOperaDataEntity.UserName); bl {
				userMap := make(map[int]int64)
				for i := 1;i <= len(userList);i++ {
					userMap[i-1] = userList[i-1]
				}
				// 返回用户数据
				c.JSON(errors.OK.HttpCode,gin.H{
					"code":errors.OK.Code,
					"message":errors.OK.Message,
					"data":userMap,
				})
			}
		}
	}
}

func (u *UserRouter) SignPost(c *gin.Context) {
	json := UserJsonBindSign{}
	_ = c.ShouldBindBodyWith(&json,binding.JSON)
	var md model.User = &model.Operations{}
	uid,username,msg_stu := md.Sign(map[string]string{
		"UserName":json.UserName,
		"UserPassWord":json.UserPassWord,
	})
	if uid == 0 {
		c.JSON(msg_stu.HttpCode,gin.H{
			"code":msg_stu.Code,
			"message":msg_stu.Message,
		})
	} else {
		var jwt auth.JwtAPI = &auth.JWT{}
		jwt.Init()
		token := jwt.NewToken(username, strconv.FormatInt(uid, 10),"")
		// 查询结果不为0则登录成功
		c.JSON(msg_stu.HttpCode,gin.H{
			"code":msg_stu.Code,
			"message":msg_stu.Message,
			"data":map[string]interface{}{
				"userdata":      uid,
				"username": username,
				"token":    token,
			},
		})
	}
}

func (u *UserRouter) UserDeleteSignState(c *gin.Context) {
	// 读取配置文件
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	if cond.ConfData.Operation.JwtStateSave == true {
		tokenHead := c.Request.Header.Get("Authorization")
		tokenHeadInfo := strings.SplitN(tokenHead," ",2)
		Claims,_ := func() (*auth.MyClaims,bool){
			var t auth.JwtAPI = &auth.JWT{}
			t.Init()
			return t.MatchToken(tokenHeadInfo[1])
		}()
		var d model.UserState = &model.OperationRedis{}
		// 删除操作成功或者失败
		if d.UserDelSignState(Claims.Uid) {
			c.JSON(errors.DelSignOK.HttpCode,gin.H{
				"code":errors.DelSignOK.Code,
				"message":errors.DelSignOK.Message,
			})
		} else {
			c.JSON(errors.ErrUserNotFound.HttpCode,gin.H{
				"code":errors.ErrUserNotFound.Code,
				"message":errors.ErrUserNotFound.Message,
			})
		}
	}
}

func (u *UserRouter) UserGetSignState(c *gin.Context) {
	// 读取配置文件
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	if cond.ConfData.Operation.JwtStateSave == true {
		uid := c.Query("uid")
		// 鉴权接口，不做验证
		var d model.UserState = &model.OperationRedis{}
		if d.UserGetSignState(uid) {
			c.JSON(errors.UserSignOk.HttpCode,gin.H{
				"code":errors.UserSignOk.Code,
				"message":errors.UserSignOk.Message,
			})
		} else {
			c.JSON(errors.ErrUserSignNotFound.HttpCode,gin.H{
				"code":errors.ErrUserSignNotFound.Code,
				"message":errors.ErrUserSignNotFound.Message,
			})
		}
	}
}



