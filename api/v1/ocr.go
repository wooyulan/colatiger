package v1

type OcrReq struct {
	ImgUrl string `form:"imgUrl" json:"imgUrl" binding:"required"`
}
