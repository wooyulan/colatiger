package repository

import (
	"colatiger/internal/models"
	"context"
	"github.com/pkg/errors"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
}

type userRepository struct {
	repo *Repository
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepository{
		repo: repo,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.repo.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
