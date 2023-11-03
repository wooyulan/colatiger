package service

import (
	v1 "colatiger/api/v1/req"
	"colatiger/internal/model"
	"colatiger/internal/pkg/ocr/strategy"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type OcrRepo interface {
	Create(ctx context.Context, record *model.Ocr) error
}

type OcrService struct {
	ocrRepo OcrRepo
}

func NewOcrService(ocrRepo OcrRepo) *OcrService {
	return &OcrService{
		ocrRepo: ocrRepo,
	}
}

func (ocr *OcrService) FindRecord() {

}

func (ocr *OcrService) OcrTextForFile(ctx *gin.Context, req v1.OcrReq) (v interface{}, err error) {
	switch req.FileType {
	// 识别身份证证件
	case "idCard":
		return ocr.IdCard(ctx, req)
	// 识别驾驶证
	default:
		return strategy.NewOcr(req.ImgUrl, &strategy.IDCard{}).Ocr()
	}
}

// IdCard 身份证识别
func (ocr *OcrService) IdCard(ctx *gin.Context, req v1.OcrReq) (v interface{}, err error) {
	data, err := strategy.NewOcr(req.ImgUrl, &strategy.IDCard{}).Ocr()
	if err != nil {
		return nil, err
	}
	result, _ := json.Marshal(data)

	var ocrRepo = &model.Ocr{
		File:     req.ImgUrl,
		FileType: req.FileType,
		Result:   string(result),
	}
	err = ocr.ocrRepo.Create(ctx, ocrRepo)
	return data, err
}
