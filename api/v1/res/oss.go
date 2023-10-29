package res

type OssFail struct {
	OriFileName string      `json:"ori_file_name"`
	Status      string      `json:"status"`
	Message     interface{} `json:"message"`
}

type OssSuccess struct {
	Key         string `json:"key"`
	PreviewUrl  string `json:"preview_url"`
	Bucket      string `json:"bucket"`
	OriFileName string `json:"ori_file_name"`
	FileSize    int64  `json:"file_size"`
}

type OssReq struct {
	Total      int          `json:"total"`
	FailNum    int          `json:"fail_num"`
	SuccessNum int          `json:"success_num"`
	OssSuccess []OssSuccess `json:"success_data"`
	OssFail    []OssFail    `json:"fail_data"`
}
