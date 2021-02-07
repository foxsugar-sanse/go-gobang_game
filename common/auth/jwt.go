package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/foxsuagr-sanse/go-gobang_game/common/config"
	"time"
)

type JwtAPI interface {
	Init()
	MatchToken(token string)  (*MyClaims,bool)
	NewToken(username string,uid string,op string) string
}

type JWT struct {
	Key []byte
	EncodeStyle string
}

type MyClaims struct {
	UserName 		string
	Uid 			string
	jwt.StandardClaims
}

type header struct {
	EncodeStyle string // 加密方式
	Type        string // Token的类型
}

type payLoad struct {
	EndTime  time.Time // 过期时间
	Username string    // 用户名
	UserId   int64     // 用户Id
}

func (j *JWT) Init() {
	// 获得token签名秘钥
	var con config.ConFig = &config.Config{}
	cond :=  con.InitConfig()
	j.Key = []byte(cond.ConfData.Jwt.Key)
}

func (j *JWT)MatchToken(tokendata string)  (*MyClaims,bool){
	// 解密token
	tokens,_ := jwt.ParseWithClaims(tokendata,&MyClaims{},func(token *jwt.Token)(interface{},error){
		return j.Key,nil
	})
	// 解密是否成功
	if key, code := tokens.Claims.(*MyClaims); code && tokens.Valid {
		return key,true
	} else {
		return nil, false
	}
}

func (j *JWT)NewToken(username string,uid string,op string) string {

	Claims := MyClaims{
		UserName:       username,
		Uid:            uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 604800,
			Issuer: op,
		},
	}
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,Claims)
	tokens,err := token.SignedString(j.Key)
	if err != nil {
		panic(err)
	} else {
		return tokens
	}
}