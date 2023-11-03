package handler

import (
	"colatiger/api/response"
	v1 "colatiger/api/v1/req"
	"colatiger/internal/pkg/img_classification"
	"colatiger/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jassue/gin-wire/app/pkg/request"
)

type OcrHandler struct {
	ocrService *service.OcrService
}

func NewOcrHandler(ocrService *service.OcrService) *OcrHandler {
	return &OcrHandler{
		ocrService: ocrService,
	}
}

func (ocr *OcrHandler) OcrTextFromFile(ctx *gin.Context) {
	var req v1.OcrReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailByErr(ctx, request.GetError(req, err))
		return
	}
	data, err := ocr.ocrService.OcrTextForFile(ctx, req)
	if err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, data)
	}
}

// ImgClassification 图片分类
func (ocr *OcrHandler) ImgClassification(ctx *gin.Context) {
	var req v1.ImgClassificationReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.FailByErr(ctx, request.GetError(req, err))
		return
	}
	res, err := img_classification.ImgClassificationPost(req)

	if err != nil {
		response.FailByErr(ctx, err)
	} else {
		response.Success(ctx, res)
	}
}
