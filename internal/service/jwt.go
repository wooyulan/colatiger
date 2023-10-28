package service

import (
	cErr "colatiger/api/error"
	"colatiger/config"
	"colatiger/internal/model"
	"colatiger/pkg/common"
	"colatiger/pkg/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtService struct {
	conf        *config.Configuration
	log         *log.Logger
	uS          *UserService
	lockBuilder *common.LockBuilder
	jRepo       JwtRepo
}

func NewJwtService(conf *config.Configuration, log *log.Logger, uS *UserService, lockBuilder *common.LockBuilder, jRepo JwtRepo) *JwtService {
	return &JwtService{
		conf:        conf,
		log:         log,
		uS:          uS,
		lockBuilder: lockBuilder,
		jRepo:       jRepo,
	}
}

type JwtRepo interface {
	JoinBlackList(ctx context.Context, tokenStr string, joinUnix int64, expires time.Duration) error
	GetBlackJoinUnix(ctx context.Context, tokenStr string) (int64, error)
}

func (s *JwtService) CreateToken(GuardName string, user model.JwtUser) (*model.TokenOutPut, *jwt.Token, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		model.CustomClaims{
			Key: GuardName,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.conf.Jwt.JwtTtl))),
				NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -1000)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        user.GetUid(),
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(s.conf.Jwt.Secret))
	if err != nil {
		return nil, nil, cErr.BadRequest("create token error:" + err.Error())
	}

	return &model.TokenOutPut{
		AccessToken: tokenStr,
		ExpiresIn:   int(s.conf.Jwt.JwtTtl),
	}, token, nil
}

func (s *JwtService) JoinBlackList(ctx *gin.Context, token *jwt.Token) error {
	nowUnix := time.Now().Unix()
	timer := token.Claims.(*model.CustomClaims).ExpiresAt.Sub(time.Now())
	fmt.Println("JoinBlackList timer", timer)

	if err := s.jRepo.JoinBlackList(ctx, token.Raw, nowUnix, timer); err != nil {
		s.log.Error(err.Error())
		return cErr.BadRequest("登出失败")
	}

	return nil
}

func (s *JwtService) IsInBlacklist(ctx *gin.Context, tokenStr string) bool {
	joinUnix, err := s.jRepo.GetBlackJoinUnix(ctx, tokenStr)
	if err != nil {
		return false
	}

	if time.Now().Unix()-joinUnix < s.conf.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true
}

func (s *JwtService) GetUserInfo(ctx *gin.Context, guardName, id string) (model.JwtUser, error) {
	switch guardName {
	case model.AppGuardName:
		return s.uS.FindUserInfoById(ctx, id)
	default:
		return nil, cErr.BadRequest("guard " + guardName + " does not exist")
	}
}

func (s *JwtService) RefreshToken(ctx *gin.Context, guardName string, token *jwt.Token) (*model.TokenOutPut, error) {
	idStr := token.Claims.(*model.CustomClaims).ID

	lock := s.lockBuilder.NewLock(ctx, "refresh_token_lock:"+idStr, s.conf.Jwt.JwtBlacklistGracePeriod)

	if lock.Get() {
		user, err := s.GetUserInfo(ctx, guardName, idStr)
		if err != nil {
			s.log.Error(err.Error())
			lock.Release()
			return nil, err
		}

		tokenData, _, err := s.CreateToken(guardName, user)
		if err != nil {
			lock.Release()
			return nil, err
		}

		err = s.JoinBlackList(ctx, token)
		if err != nil {
			lock.Release()
			return nil, err
		}

		return tokenData, nil
	}

	return nil, cErr.BadRequest("系统繁忙")
}
