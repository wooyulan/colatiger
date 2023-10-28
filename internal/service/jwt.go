package service

import (
	cErr "colatiger/api/error"
	"colatiger/internal/model"
	"colatiger/pkg/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type JwtService struct {
	conf *viper.Viper
	log  *log.Logger
	uS   *UserService
}

func NewJwtService(conf *viper.Viper, log *log.Logger, uS *UserService) *JwtService {
	return &JwtService{
		conf: conf,
		log:  log,
		uS:   uS,
	}
}

func (s *JwtService) CreateToken(GuardName string, user model.JwtUser) (*model.TokenOutPut, *jwt.Token, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		model.CustomClaims{
			Key: GuardName,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.conf.GetInt("jwt.jwt_ttl")))),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -1000)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        user.GetUid(),
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(s.conf.GetString("jwt.secret")))
	if err != nil {
		return nil, nil, cErr.BadRequest("create token error:" + err.Error())
	}

	return &model.TokenOutPut{
		AccessToken: tokenStr,
		ExpiresIn:   int(s.conf.GetInt("jwt.jwt_ttl")),
	}, token, nil
}
