package req

type OcrReq struct {
	ImgUrl   string `form:"imgUrl" json:"imgUrl" binding:"required"`
	FileType string `form:"fileType" json:"fileType" binding:"required"`
}
