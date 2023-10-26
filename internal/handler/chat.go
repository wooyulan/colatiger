package handler

import (
	"colatiger/pkg/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler interface {
	ChatStream(ctx *gin.Context)
}

func NewChatHandler(handler *Handler) ChatHandler {
	return &chatHandler{
		Handler: handler,
	}
}

type chatHandler struct {
	*Handler
}

// 对话接口
func (c chatHandler) ChatStream(ctx *gin.Context) {

	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Type", "text/event-stream; charset=utf-8")

	chat.SendMsg(ctx)

	ctx.Writer.WriteString("data: [DONE]\n\n")

}
