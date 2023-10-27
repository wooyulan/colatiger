package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/api/v1/res"
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

// 对话接口流
func (c chatHandler) ChatStream(ctx *gin.Context) {
	var form v1.ChatReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		res.ValidateFail(ctx, v1.GetErrorMsg(form, err))
		return
	}
	chat.BuildLLaVaModelBody(ctx, form)
	ctx.Writer.WriteString("data: [DONE]\n\n")

}
