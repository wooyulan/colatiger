package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/api/v1/res"
	"colatiger/internal/service"
	"colatiger/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

func NewUserHandler(handler *Handler, userService service.UserService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
	}
}

type userHandler struct {
	*Handler
	userService service.UserService
}

func (u userHandler) Register(ctx *gin.Context) {
	var form v1.Register
	if err := ctx.ShouldBindJSON(&form); err != nil {
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if err := u.userService.Register(ctx, form); err != nil {
		res.BusinessFail(ctx, err.Error())
	} else {
		res.Success(ctx, nil)
	}
}

func (u userHandler) Login(ctx *gin.Context) {
	var form v1.Login
	if err := ctx.ShouldBindJSON(&form); err != nil {
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if user, err := u.userService.Login(ctx, form); err != nil {
		res.BusinessFail(ctx, err.Error())
	} else {
		tokenData, err, _ := u.jwt.GenToken(jwt.AppGuardName, user)
		if err != nil {
			res.BusinessFail(ctx, err.Error())
			return
		}
		res.Success(ctx, tokenData)
	}
}
