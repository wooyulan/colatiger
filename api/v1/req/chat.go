package req

type ChatReq struct {
	Message      string   `form:"message" json:"message" binding:"required"`
	MaxNewTokens int      `form:"max_new_tokens" json:"max_new_tokens" `
	Temperature  float64  `form:"temperature" json:"temperature"`
	Stop         string   `form:"stop" json:"stop"`
	Images       []string `form:"images" json:"images"`
}
