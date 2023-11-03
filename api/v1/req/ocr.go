package req

// ocr 识别请求接口
type OcrReq struct {
	ImgUrl   string `form:"imgUrl" json:"imgUrl" binding:"required"`
	FileType string `form:"fileType" json:"fileType" binding:"required"`
}

// ocr 查询接口
type OcrSearchReq struct {
	Name     string `form:"name" json:"name"`
	IdCard   string `form:"idCard" json:"idCard"`
	KeyWords string `form:"keywords" json:"keywords"`
}

// ImgClassificationReq 图拍分类
type ImgClassificationReq struct {
	ImgUrl string `form:"imgUrl" json:"imgUrl" binding:"required"`
}
