package service

import (
	"colatiger/internal/repository"
	"context"
)

type UserService interface {
	Register(ctx context.Context) error
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

func (u userService) Register(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
