package repository

import (
	"colatiger/internal/model"
	"colatiger/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type chatRepository struct {
	repo *Repository
	log  *zap.Logger
}

func NewChatRepository(log *zap.Logger, repo *Repository) service.ChatRepo {
	return &chatRepository{
		log:  log,
		repo: repo,
	}
}

func (c chatRepository) SaveMessage(ctx context.Context, his *model.Chat) error {
	id, err := c.repo.sf.NextID()
	if err != nil {
		return err
	}
	his.Id = id
	if err := c.repo.db.Create(his).Error; err != nil {
		return errors.Wrap(err, "failed to create history")
	}
	return nil
}

// FindByHis 查询对话记录
func (c chatRepository) FindByHis(ctx *gin.Context) (data *[]model.Chat, err error) {
	var his []model.Chat
	if err := c.repo.db.Order("created_at desc").Limit(5).Find(&his).Error; err != nil {
		fmt.Printf(err.Error())
		return nil, errors.Wrap(err, "查询失败")
	}
	return &his, nil
}

func (c chatRepository) DelChatHisByChatIdAndUserId(assistant string, user string) {
	c.repo.db.Delete(model.Chat{})
}
