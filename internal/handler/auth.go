package handler

import (
	"colatiger/api/response"
	v1 "colatiger/api/v1"
	"colatiger/internal/model"
	"colatiger/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log         *zap.Logger
	jwtService  *service.JwtService
	userService *service.UserService
}

func NewAuthHandler(log *zap.Logger, jwtService *service.JwtService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (u *AuthHandler) Register(ctx *gin.Context) {
	var form v1.Register
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.FailByErr(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if err := u.userService.Register(ctx, form); err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

func (u *AuthHandler) Login(ctx *gin.Context) {
	var form v1.Login
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.FailByErr(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if user, err := u.userService.Login(ctx, form); err != nil {
		response.FailByErr(ctx, err)
	} else {
		tokenData, _, err := u.jwtService.CreateToken(model.AppGuardName, user)
		if err != nil {
			response.FailByErr(ctx, err)
			return
		}
		response.Success(ctx, tokenData)
	}
}

func (u *AuthHandler) GetInfo(ctx *gin.Context) {
	user, err := u.userService.FindUserInfoById(ctx, ctx.Keys["id"].(string))
	if err != nil {
		response.FailByErr(ctx, err)
		return
	}
	response.Success(ctx, user)
}
