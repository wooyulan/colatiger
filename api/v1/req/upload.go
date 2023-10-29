package req

import "mime/multipart"

type UploadReq struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}
