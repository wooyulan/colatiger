package handler

import (
	v1 "colatiger/api/v1"
	"colatiger/api/v1/res"
	"colatiger/pkg/chat"
	"colatiger/pkg/log"
	"github.com/gin-gonic/gin"
)

func NewChatHandler(log *log.Logger) *ChatHandler {
	return &ChatHandler{
		log: log,
	}
}

type ChatHandler struct {
	log *log.Logger
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
