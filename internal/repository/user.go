package repository

import (
	"colatiger/internal/models"
	"colatiger/internal/service"
	"colatiger/pkg/log"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	repo *Repository
	log  *log.Logger
}

func NewUserRepository(log *log.Logger, repo *Repository) service.UserRepo {
	return &userRepository{
		log:  log,
		repo: repo,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.repo.db.Create(user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

// FindByEmail 根据邮箱查询
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID 根据id主键查询
func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.repo.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.repo.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
