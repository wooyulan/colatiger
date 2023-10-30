package handler

import (
	"colatiger/api/response"
	"colatiger/api/v1/req"
	"colatiger/internal/service"
	"colatiger/pkg/chat"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewChatHandler(log *zap.Logger, milvusS *service.MilvusService) *ChatHandler {
	return &ChatHandler{
		log:     log,
		milvusS: milvusS,
	}
}

type ChatHandler struct {
	log     *zap.Logger
	milvusS *service.MilvusService
}

// 对话接口流
func (c *ChatHandler) ChatStream(ctx *gin.Context) {
	var form req.ChatReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.FailByErr(ctx, req.GetErrorMsg(form, err))
		return
	}
	chat.BuildLLaVaModelBody(ctx, form)
	ctx.Writer.WriteString("data: [DONE]\n\n")

}

func (c *ChatHandler) Test(ctx *gin.Context) {
	ok, err := c.milvusS.InsertData()
	if err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, ok)
	}

}
