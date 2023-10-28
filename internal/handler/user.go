package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/api/v1/res"
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
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if err := u.userService.Register(ctx, form); err != nil {
		res.BusinessFail(ctx, "")
	} else {
		res.Success(ctx, nil)
	}
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var form v1.Login
	if err := ctx.ShouldBindJSON(&form); err != nil {
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}

	if user, err := u.userService.Login(ctx, form); err != nil {
		res.BusinessFail(ctx, err.Error())
	} else {
		tokenData, err := u.jwt.GenToken(strconv.FormatInt(user.Id, 10), time.Now().Add(time.Hour*24*90))
		if err != nil {
			res.BusinessFail(ctx, err.Error())
			return
		}
		res.Success(ctx, tokenData)
	}
}

func (u *UserHandler) GetInfo(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		res.TokenFail(ctx)
		return
	}
	user, err := u.userService.FindUserInfoById(ctx, userId)
	if err != nil {
		res.TokenFail(ctx)
		return
	}

	res.Success(ctx, user)
}
