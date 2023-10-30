package service

import "github.com/pkg/errors"

type MilvusRepo interface {
	HasCollection() (bool, error)
}

type MilvusService struct {
	repo MilvusRepo
}

func NewMilvusService(repo MilvusRepo) *MilvusService {
	return &MilvusService{
		repo: repo,
	}
}

func (m *MilvusService) InsertData() (ok bool, err error) {
	if ok, err := m.repo.HasCollection(); err != nil || !ok {
		return ok, errors.New("Collection 不存在")
	}
	return
}
