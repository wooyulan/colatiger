package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"regexp"
	"time"
)

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

// 所有需要颁发 token 的用户模型必须实现这个接口
type JwtUser interface {
	GetUid() string
}

type JWT struct {
	key                []byte
	ttl                time.Time `mapstructure:"jwt_ttl" json:"ttl" yaml:"ttl"`                                                // token 有效期（秒）
	refreshGracePeriod int64     `mapstructure:"refresh_grace_period" json:"refresh_grace_period" yaml:"refresh_grace_period"` // token 自动刷新宽限时间（秒）

}

// MyCustomClaims 自定义Claims
type MyCustomClaims struct {
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	return &JWT{key: []byte(conf.GetString("security.jwt.key"))}
}

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (j *JWT) GenToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error, token *jwt.Token) {
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(j.ttl),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    GuardName, // 用于在中间件中区分不同客户端颁发的 token，避免 token 跨端使用,
			ID:        user.GetUid(),
		},
	})

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	tokenData = TokenOutPut{
		tokenString,
		TokenType,
	}
	return
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	re := regexp.MustCompile(`(?i)Bearer `)
	tokenString = re.ReplaceAllString(tokenString, "")
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
