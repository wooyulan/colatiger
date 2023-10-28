package handler

import (
	"colatiger/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}

var ProviderSet = wire.NewSet(NewUserHandler, NewChatHandler)
