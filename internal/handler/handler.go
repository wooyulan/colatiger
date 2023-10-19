package handler

import (
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
	jwt    *jwt.JWT
}

func NewHandler(logger *log.Logger, jwt *jwt.JWT) *Handler {
	return &Handler{
		logger: logger,
		jwt:    jwt,
	}
}

func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
