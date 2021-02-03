package errors

import (
	"net/http"
)

/*
错误码设计
第一位表示错误级别, 1 为系统错误, 2 为普通错误
第二三四位表示服务模块代码
第五六位表示具体错误代码
*/

type Errno struct {
	Code     int    // 错误码
	Message  string // 展示给用户看的
	HttpCode int    // HTTP状态码
}

var (
	OK = &Errno{Code: 0, Message: "OK", HttpCode: http.StatusOK}

	// 系统错误
	ErrUnKnown        = &Errno{Code: 100000, Message: "未知错误", HttpCode: http.StatusInternalServerError}
	ErrInternalServer = &Errno{Code: 100001, Message: "内部服务器错误", HttpCode: http.StatusInternalServerError}
	ErrParamConvert   = &Errno{Code: 100002, Message: "参数转换时发生错误", HttpCode: http.StatusInternalServerError}
	ErrDatabase       = &Errno{Code: 100003, Message: "数据库错误", HttpCode: http.StatusInternalServerError}

	// 模块通用错误
	ErrValidation      = &Errno{Code: 200001, Message: "参数校验失败", HttpCode: http.StatusForbidden}
	ErrBadRequest      = &Errno{Code: 200002, Message: "请求参数错误", HttpCode: http.StatusBadRequest}
	ErrGetTokenFail    = &Errno{Code: 200003, Message: "获取 token 失败", HttpCode: http.StatusForbidden}
	ErrTokenNotFound   = &Errno{Code: 200004, Message: "用户 token 不存在", HttpCode: http.StatusUnauthorized}
	ErrTokenExpire     = &Errno{Code: 200005, Message: "用户 token 过期", HttpCode: http.StatusForbidden}
	ErrTokenValidation = &Errno{Code: 200005, Message: "用户 token 无效", HttpCode: http.StatusForbidden}
	ErrTimeNoSwitch	   = &Errno{Code: 200006, Message: "提交时间不正确",HttpCode: http.StatusForbidden}

	// User模块错误
	ErrUserNotFound       = &Errno{Code: 200104, Message: "用户不存在", HttpCode: http.StatusBadRequest}
	ErrPasswordIncorrect  = &Errno{Code: 200105, Message: "密码错误", HttpCode: http.StatusBadRequest}
	ErrUserRegisterAgain  = &Errno{Code: 200107, Message: "重复注册", HttpCode: http.StatusBadRequest}
	ErrUserNameExist	  = &Errno{Code: 200108, Message: "用户名已存在", HttpCode: http.StatusBadRequest}
	ErrUsernameValidation = &Errno{Code: 200109, Message: "用户名不合法", HttpCode: http.StatusBadRequest}
	ErrPasswordValidation = &Errno{Code: 200110, Message: "密码不合法", HttpCode: http.StatusBadRequest}

	// Group模块错误
)

func GetErrorsStruct(failed *Errno) *Errno {
	return failed
}