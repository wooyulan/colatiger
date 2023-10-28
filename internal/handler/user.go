package handler

import (
	"colatiger/api/response"
	v1 "colatiger/api/v1"
	"colatiger/internal/service"
	"colatiger/pkg/jwt"
	"colatiger/pkg/log"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type UserHandler struct {
	log         *log.Logger
	jwt         *jwt.JWT
	userService *service.UserService
}

func NewUserHandler(log *log.Logger, userService *service.UserService) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userService,
	}
}

func (u *UserHandler) Register(ctx *gin.Context) {
	var form v1.Register
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if err := u.userService.Register(ctx, form); err != nil {
		response.BusinessFail(ctx, "")
	} else {
		response.Success(ctx, nil)
	}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var form v1.Login
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if user, err := u.userService.Login(ctx, form); err != nil {
		response.BusinessFail(ctx, err.Error())
	} else {
		tokenData, err := u.jwt.GenToken(strconv.FormatInt(user.Id, 10), time.Now().Add(time.Hour*24*90))
		if err != nil {
			response.BusinessFail(ctx, err.Error())
			return
		}
		response.Success(ctx, tokenData)
	}
}

func (u *UserHandler) GetInfo(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		response.TokenFail(ctx)
		return
	}
	user, err := u.userService.FindUserInfoById(ctx, userId)
	if err != nil {
		response.TokenFail(ctx)
		return
	}

	response.Success(ctx, user)
}
