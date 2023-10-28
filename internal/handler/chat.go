package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/api/v1/res"
	"colatiger/pkg/chat"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IChatHandler interface {
	ChatStream(ctx *gin.Context)
}

func NewChatHandler(log *zap.Logger) *ChatHandler {
	return &ChatHandler{
		log: log,
	}
}

type ChatHandler struct {
	log *zap.Logger
}

// 对话接口流
func (c *ChatHandler) ChatStream(ctx *gin.Context) {
	var form v1.ChatReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}
	chat.BuildLLaVaModelBody(ctx, form)
	ctx.Writer.WriteString("data: [DONE]\n\n")

}
