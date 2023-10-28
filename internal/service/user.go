package service

import (
	v1 "colatiger/api/v1"
	"colatiger/internal/models"
	"colatiger/internal/repository"
	"colatiger/pkg/helper/hash"
	"colatiger/pkg/helper/sid"
	"context"
	"github.com/pkg/errors"
)

type IUserService interface {
	Register(ctx context.Context, reg v1.Register) error
	Login(ctx context.Context, login v1.Login) (user *models.User, err error)
	FindUserInfoById(ctx context.Context, userId string) (user *models.User, err error)
}

type UserService struct {
	userRepo repository.UserRepository
	sid      *sid.Sid
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (u *UserService) Login(ctx context.Context, login v1.Login) (user *models.User, err error) {
	user, err = u.userRepo.FindByEmail(ctx, login.Username)
	if err != nil || !hash.BcryptMakeCheck([]byte(login.Password), user.Password) {
		err = errors.Wrap(err, "用户不存在或密码错误")
	}
	return
}

// Register 用户注册
func (u *UserService) Register(ctx context.Context, req v1.Register) error {
	if user, err := u.userRepo.FindByEmail(ctx, req.Email); err == nil && user != nil {
		return errors.New("用户名已经存在")
	}
	primaryKey, err := u.sid.GenInt64()
	if err != nil {
		return errors.Wrap(err, "failed to generate user ID")
	}
	// Create a user
	var user = &models.User{
		Id:       primaryKey,
		Password: hash.BcryptMake([]byte(req.Password)),
		Mobile:   req.Mobile,
		Email:    req.Email,
	}
	if err = u.userRepo.Create(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

// 获取根据用户id获取
func (u *UserService) FindUserInfoById(ctx context.Context, userId string) (user *models.User, err error) {
	if user, err = u.userRepo.FindByID(ctx, userId); err != nil || user == nil {
		return nil, errors.New("当前用户不存在")
	}
	user.Password = ""
	return
}
