package handler

import "github.com/gin-gonic/gin"

type OcrHandler interface {
	OcrImageInfo(ctx *gin.Context)
}

type ocrHandler struct {
	*Handler
}

func NewOcrHandler(handler *Handler) OcrHandler {
	return &ocrHandler{
		Handler: handler,
	}
}

func (o ocrHandler) OcrImageInfo(ctx *gin.Context) {

	//TODO implement me
	panic("implement me")
}
