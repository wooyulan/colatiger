package service

import (
	v1 "colatiger/api/v1/req"
	"colatiger/internal/model"
	"colatiger/pkg/chat"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

type ChatRepo interface {
	SaveMessage(ctx context.Context, his *model.Chat) error
	FindByHis(ctx *gin.Context) (data *[]model.Chat, err error)
}

type ChatService struct {
	chatRepo ChatRepo
}

func NewChatService(repo ChatRepo) *ChatService {
	return &ChatService{
		chatRepo: repo,
	}
}

func (c *ChatService) SyncChatMessage(ctx *gin.Context, req v1.ChatReq) error {
	// 构建对话历史记录
	data, _ := c.findHistory(ctx)

	assistantRes := chat.BuildLLaVaModelBody(ctx, req, data)

	ctx.Writer.WriteString("data: [DONE]\n\n")
	// 存储聊天记录
	/**
	todo 用户id  模型参数
	*/
	var chatHis = &model.Chat{
		Question: req.Message,
		Answer:   assistantRes,
		Prompt:   req.Prompt,
	}
	if req.Images != nil && len(req.Images) > 0 {
		chatHis.File = strings.Join(req.Images, ",")
	}
	if err := c.chatRepo.SaveMessage(ctx, chatHis); err != nil {
		return errors.Wrap(err, "failed to create hitory")
	}
	return nil
}

// 查询对话聊天记录
func (c *ChatService) findHistory(ctx *gin.Context) (*[]model.Chat, error) {
	data, err := c.chatRepo.FindByHis(ctx)
	return data, err

}
