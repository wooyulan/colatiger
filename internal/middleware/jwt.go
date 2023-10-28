package middleware

import (
	cErr "colatiger/api/error"
	"colatiger/api/response"
	"colatiger/config"
	"colatiger/internal/model"
	"colatiger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type JWTAuth struct {
	conf *config.Configuration
	jwtS *service.JwtService
}

func NewJWTAuth(
	conf *config.Configuration,
	jwtS *service.JwtService,
) *JWTAuth {
	return &JWTAuth{
		conf: conf,
		jwtS: jwtS,
	}
}

func (m *JWTAuth) Handler(guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.FailByErr(c, cErr.Unauthorized("missing Authorization header"))
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.conf.Jwt.Secret), nil
		})
		if err != nil || m.jwtS.IsInBlacklist(c, tokenStr) {
			response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
			return
		}

		claims := token.Claims.(*model.CustomClaims)
		if claims.Issuer != guardName {
			response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
			return
		}

		// token 续签
		if int64(claims.ExpiresAt.Sub(time.Now()).Seconds()) < m.conf.Jwt.RefreshGracePeriod {
			tokenData, err := m.jwtS.RefreshToken(c, guardName, token)
			if err == nil {
				c.Header("new-token", tokenData.AccessToken)
				c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
			}
		}

		c.Set("token", token)
		c.Set("id", claims.ID)
	}
}
