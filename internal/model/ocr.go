package model

type Ocr struct {
	Id       uint64 `json:"id" gorm:"primaryKey"`
	File     string `json:"file" gorm:"comment:识别文件"`
	Result   string `json:"result" gorm:"comment:识别结果"`
	UserId   uint64 `json:"user_id" gorm:"comment:所属用户"`
	FileType string `json:"file_type" gorm:"comment:文件类型"`
	Timestamps
	SoftDeletes
}

func (u *Ocr) TableName() string {
	return "t_ocr_record"
}
