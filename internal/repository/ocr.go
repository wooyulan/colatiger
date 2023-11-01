package repository

import (
	"colatiger/internal/model"
	"colatiger/internal/service"
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ocrRepository struct {
	repo *Repository
	log  *zap.Logger
}

func NewOcrRepository(log *zap.Logger, repo *Repository) service.OcrRepo {
	return &ocrRepository{
		log:  log,
		repo: repo,
	}
}

func (o ocrRepository) Create(ctx context.Context, record *model.Ocr) error {
	id, err := o.repo.sf.NextID()
	if err != nil {
		return err
	}
	record.Id = id
	if err := o.repo.db.Create(record).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
