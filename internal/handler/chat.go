package handler

import (
	"colatiger/api/response"
	"colatiger/api/v1/req"
	"colatiger/internal/service"
	"colatiger/pkg/helper/convert"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewChatHandler(log *zap.Logger, milvusS *service.MilvusService, chatS *service.ChatService) *ChatHandler {
	return &ChatHandler{
		log:     log,
		milvusS: milvusS,
		chatS:   chatS,
	}
}

type ChatHandler struct {
	log     *zap.Logger
	milvusS *service.MilvusService
	chatS   *service.ChatService
}

// 对话接口流
func (c *ChatHandler) ChatStream(ctx *gin.Context) {
	var form req.ChatReq
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.FailByErr(ctx, req.GetErrorMsg(form, err))
		return
	}

	form.Images = convert.SliceRemoveStr(form.Images, "")

	c.chatS.SyncChatMessage(ctx, form)
}

// 查询对话记录
func (c *ChatHandler) FindChatHis(ctx *gin.Context) {
	if data, err := c.chatS.FindHistory(ctx); err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, data)
	}
}

// 删除 对话记录
func (c *ChatHandler) DelChatHis(ctx *gin.Context) {
	//userId := ctx.Keys["id"].(string
	userId := ""
	assistant := ""
	c.chatS.DelChatHisByChatAndUserId(assistant, userId)
	response.Success(ctx, nil)
}

// Test 测试向量库
func (c *ChatHandler) Test(ctx *gin.Context) {
	ok, err := c.milvusS.InsertData()
	if err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, ok)
	}
}
