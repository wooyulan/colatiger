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
	Register(ctx context.Context, register v1.Register) error
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

// Register 用户注册
func (u userService) Register(ctx context.Context, req v1.Register) error {
	// Generate user ID
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
	}
	if err = u.userRepo.Create(ctx, user); err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
