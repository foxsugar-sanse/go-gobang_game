package controller

import (
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	//"time"
	//"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
)

type RouterRequest interface {
	LoginPost			(c * gin.Context)
	UserInfoGet			(c * gin.Context)
	UserInfoUpdate		(c * gin.Context)
	UserDelete			(c * gin.Context)
	UserOtherOperations	(c * gin.Context)
	SignPost			(c * gin.Context)
}

type UserJsonBindLogin struct {
	TimeSlices 		int64  `json:"t_s"`
	UserName 		string `json:"user_name"`
	UserPassword 	string `json:"user_password"`
}

type UserJsonBindSign struct {
	UserName		string `json:"user_name"`
	UserPassWord	string `json:"user_pass_word"`
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
	json := &UserJsonBindLogin{}
	_ = c.BindJSON(&UserJsonBindLogin{})

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
	// 鉴权通过，获取用户token，得知id返回对应信息
	tokenHead := c.Request.Header.Get("Authorization")
	tokenHeadInfo := strings.SplitN(tokenHead," ",2)
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	MyClaims,Op := jwt.MatchToken(tokenHeadInfo[1])
	if Op == true {
		// token解密成功
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(MyClaims.Uid,10,64)
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
	panic("implement me")
}

func (u *UserRouter) SignPost(c *gin.Context) {
	json := &UserJsonBindSign{}
	_ = c.BindJSON(&UserJsonBindSign{})
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




