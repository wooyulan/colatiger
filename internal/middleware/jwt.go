package middleware

import (
	"colatiger/api/v1/res"
	"colatiger/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func StrictAuth(j *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			res.TokenFail(ctx)
			ctx.Abort()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			res.TokenFail(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func NoStrictAuth(j *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
