package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(ctx *gin.Context)
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
		v1.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if err := u.userService.Register(ctx, form); err != nil {
		v1.BusinessFail(ctx, err.Error())
	} else {
		v1.Success(ctx, nil)
	}
}
