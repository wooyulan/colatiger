package repository

import (
	"colatiger/internal/service"
	"context"
	"go.uber.org/zap"
)

type milvusRepository struct {
	repo *Repository
	log  *zap.Logger
}

func NewMilvusRepository(repo *Repository, log *zap.Logger) service.MilvusRepo {
	return &milvusRepository{
		repo: repo,
		log:  log,
	}
}

func (m *milvusRepository) HasCollection() (bool, error) {
	has, err := m.repo.milvus.HasCollection(context.Background(), "tenant")
	return has, err
}
