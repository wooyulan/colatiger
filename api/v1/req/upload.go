package v1

import "mime/multipart"

type UploadReq struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}
