package controller

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/foxsuagr-sanse/go-gobang_game/app/model"
	"github.com/foxsuagr-sanse/go-gobang_game/common/auth"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"github.com/foxsuagr-sanse/go-gobang_game/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/tencentyun/cos-go-sdk-v5"
)


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
	FormGroupGetUserFriend  (c * gin.Context)
	AddUserForFriend		(c * gin.Context)
	DeleteUserForFriend		(c * gin.Context)
	ModifyFriendInfo		(c * gin.Context)
	OtherUserFriendInterface(c * gin.Context)

	RefuseUserFriendRequest	(c * gin.Context)
	ConsentUserFriendRequest(c * gin.Context)
	GetUserFriendRequest	(c * gin.Context)
	CreateUserFriendRequest	(c * gin.Context)

	CreateUserPortrait		(c * gin.Context)
	DeleteUserPortrait	    (c * gin.Context)
	// 和控制用户分组的接口组合
	UserGroupController
	// 和控制用户游戏邀请的接口结合
	UserGames
}

type UserRouter struct {
	OperationUserGroup
	UserInviteFunc
	//UserJBL *UserJsonBindLogin
	//UserJBS *UserJsonBindSign
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
	UserAge			int    `json:"user_age"`
	UserContact		string `json:"user_contact"`
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
	idArgs:= c.Query("uid")
	if idArgs == "self" {
		// ?self 获取关于自己的详细信息
		var jwt auth.JwtAPI = &auth.JWT{}
		jwt.Init()
		tokenHead := c.Request.Header.Get("Authorization")
		tokenHeadInfo := strings.SplitN(tokenHead, " ", 2)
		if token,bl := jwt.MatchToken(tokenHeadInfo[1]); bl {
			uid,_ := strconv.ParseInt(token.Uid,10,64)
			users := md.GetUserInfo(uid)
			c.JSON(errors.OK.HttpCode,gin.H{
				"code":errors.OK.Code,
				"message":errors.OK.Message,
				"data":map[string]interface{}{
					"uid":             users.Uid,
					"user_name":       users.UserName,
					"user_nick_name":  users.UserNickName,
					"user_brief":      users.UserBrief,
					"user_age":        users.UserAge,
					"user_sex":        users.UserSex,
					"user_contact":    users.UserContact,
					"user_portrait":   users.UserPortrait,
					"user_email":	   users.UserEmail,
					"user_phone":      users.UserPhone,
					"email_validation":users.EmailValidation,
					"phone_validation":users.PhoneValidation,
				},
			})
		}
	} else {
		uid,_ := strconv.ParseInt(idArgs,10,64)
		users := md.GetUserInfo(uid)
		c.JSON(errors.OK.HttpCode,gin.H{
			"code":errors.OK.Code,
			"message":errors.OK.Message,
			"data":map[string]interface{}{
				"uid":            users.Uid,
				"user_name":      users.UserName,
				"user_nick_name": users.UserNickName,
				"user_age":       users.UserAge,
				"user_sex":       users.UserSex,
				"user_contact":   users.UserContact,
				"user_portrait":  users.UserPortrait,
			},
		})
	}
}

func (u *UserRouter) UserInfoUpdate(c *gin.Context) {
	// 绑定提交json
	json := UserJsonBindUpdate{}
	_ = c.BindJSON(&json)
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
	json := UserBindJsonOtherOpera{}
	_ = c.BindJSON(&json)

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
					userMap[i] = userList[i-1]
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
	uid,username, msgStu := md.Sign(map[string]string{
		"UserName":json.UserName,
		"UserPassWord":json.UserPassWord,
	})
	if uid == 0 {
		c.JSON(msgStu.HttpCode,gin.H{
			"code":    msgStu.Code,
			"message": msgStu.Message,
		})
	} else {
		var jwt auth.JwtAPI = &auth.JWT{}
		jwt.Init()
		token := jwt.NewToken(username, strconv.FormatInt(uid, 10),"")
		// 查询结果不为0则登录成功
		// 向redis写入登录状态
		var md model.UserState = &model.OperationRedis{}
		if md.UserCreateSignState(strconv.FormatInt(uid, 10)) {
			c.JSON(msgStu.HttpCode, gin.H{
				"code":    msgStu.Code,
				"message": msgStu.Message,
				"data": map[string]interface{}{
					"userdata": uid,
					"username": username,
					"token":    token,
				},
			})
		} else {

		}
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

func (u *UserRouter) CreateUserPortrait(c *gin.Context) {
	// TODO:用户头像上传模块
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(errors.ErrUserUploadNotFound.HttpCode,gin.H{
			"code": errors.ErrUserUploadNotFound.Code,
			"message": errors.ErrUserUploadNotFound.Message,
		})
	} else {
		// 判断文件名是否满足要求
		// 载入配置
		var con config.ConFig = &config.Config{}
		conf := con.InitConfig()
		filenames := strings.Split(conf.ConfData.Model.Contentfilename,",")
		fileformat := strings.Split(file.Filename,".")
		for i := 0; i < len(filenames); i++ {
			if fileformat[len(fileformat) - 1] == filenames[i] {
				// 判断文件大小
				if file.Size / (1024 * 1024) > int64(conf.ConfData.Model.UploadMax) {
					c.JSON(errors.ErrUserUploadExceedMax.HttpCode,gin.H{
						"code":errors.ErrUserUploadExceedMax.Code,
						"message":errors.ErrUserUploadExceedMax.Message,
					})
					break
				}
				var jwt auth.JwtAPI = &auth.JWT{}
				jwt.Init()
				tokenHead := c.Request.Header.Get("Authorization")
				tokenInfo := strings.SplitN(tokenHead, " ", 2)
				if token,err2 := jwt.MatchToken(tokenInfo[1]); err2 == true {
					if conf.ConfData.Model.Imgsave == "local" {
						// 更改文件名保存文件{本地存储}
						filename := sha256.New()
						filename.Write([]byte(token.Uid + file.Filename))
						file.Filename = hex.EncodeToString(filename.Sum(nil))
						// 先打开数据库把url存储,格式{xxx.jpg}
						var md model.User = &model.Operations{}
						uid,_ := strconv.ParseInt(token.Uid,10,64)
						if errno := md.SetUserPortraitUrl(uid,file.Filename + "." + fileformat[len(fileformat) - 1]); errno.HttpCode == 200 {
							if err3 := c.SaveUploadedFile(file,conf.ConfData.Model.Localurl + "/" + file.Filename + "." + fileformat[len(fileformat) - 1]); err3 == nil {
								c.JSON(errno.HttpCode,gin.H{
									"code":errno.Code,
									"message":errno.Message,
								})
							} else {
								c.JSON(errors.ErrUserUploadUrlNot.HttpCode,gin.H{
									"code":errors.ErrUserUploadUrlNot.Code,
									"message":errors.ErrUserUploadUrlNot.Message,
								})
							}
						} else {
							c.JSON(errno.HttpCode,gin.H{
								"code":errno.Code,
								"message":errno.Message,
							})
						}

					} else if conf.ConfData.Model.Imgsave == "tencentcloud" {
						// 在腾讯COS上存储
						filename := sha256.New()
						filename.Write([]byte(token.Uid + file.Filename))
						file.Filename = hex.EncodeToString(filename.Sum(nil))
						file.Filename = file.Filename + "." + fileformat[len(fileformat) - 1]
						fileSrc,err  := file.Open()
						if err != nil {
							panic(err)
						}
						buf := bytes.NewBuffer(nil)
						_, _ = io.Copy(buf, fileSrc)
						bufReader := strings.NewReader(buf.String())

						// 使用永久密钥连接腾讯COS
						// 桶的路径
						u,_ := url.Parse(conf.ConfData.Tencentcloud.Bucketurl)
						urlSplit := strings.SplitN(conf.ConfData.Tencentcloud.Bucketurl,".",5)
						CosRegion := urlSplit[2]
						// 地区url
						su,_ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", CosRegion))
						b := cos.BaseURL{BucketURL: u,ServiceURL: su}

						// 开始连接
						client := cos.NewClient(&b,&http.Client{
							Transport: &cos.AuthorizationTransport{
								SecretID:  conf.ConfData.Tencentcloud.Secretid,
								SecretKey: conf.ConfData.Tencentcloud.Secretkey,
							},
						})
						if client != nil {
							// 调用cos请求
							// 上传文件到存储桶
							_,err2 := client.Object.Put(context.Background(),file.Filename,bufReader,nil)
							if err2 != nil {
								c.JSON(errors.ErrTencentCosUploadNot.HttpCode,gin.H{
									"code":errors.ErrTencentCosUploadNot.Code,
									"message":errors.ErrTencentCosUploadNot.Message,
								})
							} else {
								// 将Url存入数据库
								var md model.User = &model.Operations{}
								uid,_ := strconv.ParseInt(token.Uid,10,64)
								errno := md.SetUserPortraitUrl(uid,conf.ConfData.Tencentcloud.Bucketurl + "/" + file.Filename)
								c.JSON(errno.HttpCode,gin.H{
									"code":errno.Code,
									"message":errno.Message,
								})
							}
						} else {
							c.JSON(errors.ErrTencentCosLinkError.HttpCode,gin.H{
								"code":errors.ErrTencentCosLinkError.Code,
								"message":errors.ErrTencentCosLinkError.Message,
							})
						}
					}
				}
				break
			} else if i + 1 == len(filenames) {
				// 匹配了全部条件还未匹配到则返回错误
				c.JSON(errors.ErrUserUploadFormatNo.HttpCode,gin.H{
					"code":errors.ErrUserUploadFormatNo.Code,
					"message":errors.ErrUserUploadFormatNo.Message,
				})
			}
		}

	}
}

func (u *UserRouter) DeleteUserPortrait(c *gin.Context) {
	// 删除头像
	var jwt auth.JwtAPI = &auth.JWT{}
	jwt.Init()
	tokenHead := c.Request.Header.Get("Authorization")
	tokenInfo := strings.SplitN(tokenHead, " ", 2)
	if token,err2 := jwt.MatchToken(tokenInfo[1]); err2 == true {
		var md model.User = &model.Operations{}
		uid,_ := strconv.ParseInt(token.Uid,10,64)
		errno := md.DeleteUserPortrait(uid)
		c.JSON(errno.HttpCode,gin.H{
			"code":errno.Code,
			"message":errno.Message,
		})
	}
}


