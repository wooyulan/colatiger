package service

import (
	v1 "colatiger/api/v1"
	"colatiger/internal/models"
	"colatiger/internal/repository"
	"colatiger/pkg/helper/hash"
	"context"
	"github.com/pkg/errors"
)

type UserService interface {
	Register(ctx context.Context, reg v1.Register) error
	Login(ctx context.Context, login v1.Login) (user *models.User, err error)
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

// Login 用户登录
func (u userService) Login(ctx context.Context, login v1.Login) (user *models.User, err error) {
	user, err = u.userRepo.FindByEmail(ctx, login.Username)
	if err != nil || !hash.BcryptMakeCheck([]byte(login.Password), user.Password) {
		err = errors.Wrap(err, "用户不存在或密码错误")
	}
	return
}

// Register 用户注册
func (u userService) Register(ctx context.Context, req v1.Register) error {
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
		Name:     req.Name,
		Password: hash.BcryptMake([]byte(req.Password)),
		Mobile:   req.Mobile,
		Email:    req.Email,
	}
	if err = u.userRepo.Create(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
