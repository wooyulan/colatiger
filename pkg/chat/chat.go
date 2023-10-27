package chat

import (
	"bufio"
	"bytes"
	v1 "colatiger/api/v1"
	"colatiger/pkg/helper/img"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

type LLaVaChatReq struct {
	Model        string   `json:"model"`
	Prompt       string   `json:"prompt"`
	MaxNewTokens int      `json:"max_new_tokens"`
	Temperature  float64  `json:"temperature"`
	Stop         string   `json:"stop"`
	Images       []string `json:"images"`
}

type MessageResponse struct {
	Text      string `json:"text"`
	ErrorCode int    `json:"error_code"`
}

func SendMsg(ctx *gin.Context, llavaReq LLaVaChatReq) {

	jsonBytes, err := json.Marshal(llavaReq)
	if err != nil {
	}
	//
	req, err := http.NewRequest("POST", "http://82.156.138.158:10000/worker_generate_stream", bytes.NewBuffer(jsonBytes))
	if err != nil {
		println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	for {
		reader := bufio.NewReader(resp.Body)
		rawLine, readErr := reader.ReadBytes(0)
		if readErr != nil {
			break
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		// 删除最后一个分隔符
		noZeroLine := noSpaceLine[0 : len(noSpaceLine)-1]

		jsonStream := string(noZeroLine)
		dec := json.NewDecoder(strings.NewReader(jsonStream))
		var m MessageResponse
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			break
		}

		if m.ErrorCode == 0 {
			text := strings.Split(m.Text, llavaReq.Prompt)[1]
			llavaReq.Prompt = llavaReq.Prompt + text
			text = "\"" + text + "\""
			var data = "data: " + text + "\n\n"
			fmt.Fprintf(ctx.Writer, data)
			ctx.Writer.Flush()
		}
	}
}

func BuildLLaVaModelBody(ctx *gin.Context, chatReq v1.ChatReq) {
	prompt := "A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions. USER: %s ASSISTANT:"
	body := LLaVaChatReq{
		Model:        "llava-v1.5-13b",
		Prompt:       prompt,
		MaxNewTokens: 512,
		Temperature:  0.7,
		Stop:         "</s>",
		Images:       []string{},
	}

	if chatReq.Images != nil && len(chatReq.Images) > 0 {
		baseImg := make([]string, len(chatReq.Images))
		// base64 图片
		for i, imgUrl := range chatReq.Images {
			base64, _ := img.GetUrlImgBase64(imgUrl)
			baseImg[i] = base64
		}
		body.Images = baseImg
		body.Prompt = fmt.Sprintf(prompt, "<image>\\n"+chatReq.Message)
	} else {
		body.Prompt = fmt.Sprintf(prompt, chatReq.Message)
	}

	SendMsg(ctx, body)
}
