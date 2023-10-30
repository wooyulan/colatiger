package handler

import (
	"colatiger/pkg/third"
	"github.com/gin-gonic/gin"
)

type ThirdHandler struct {
	ocr *third.OCR
}

func NewThirdHandler(ocr *third.OCR) *ThirdHandler {
	return &ThirdHandler{
		ocr: ocr,
	}
}

func (t *ThirdHandler) IdCard(ctx *gin.Context) {

}
