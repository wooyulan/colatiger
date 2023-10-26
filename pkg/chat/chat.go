package chat

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

type ChatReq struct {
	Model        string   `json:"model"`
	Prompt       string   `json:"prompt"`
	MaxNewTokens int      `json:"max_new_tokens"`
	Temperature  float64  `json:"temperature"`
	Stop         string   `json:"stop"`
	Images       []string `json:"images"`
}

type Message struct {
	Text      string `json:"text"`
	ErrorCode int    `json:"error_code"`
}

func SendMsg(ctx *gin.Context) {
	// func main() {

	body := ChatReq{
		Model:        "llava-v1.5-13b",
		Prompt:       "A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions. USER: 你好,推荐几本中国历史书籍 ASSISTANT:",
		MaxNewTokens: 512,
		Temperature:  0.7,
		Stop:         "</s>",
		Images:       []string{},
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		// handle error
	}
	//
	req, err := http.NewRequest("POST", "http://82.156.138.158:10000/worker_generate_stream", bytes.NewBuffer(jsonBytes))
	if err != nil {
		println(err)
	}
	req.Header.Set("Content-Type", "text/event-stream; charset=UTF-8")
	req.Header.Set("User-Agent", "LLaVA Client")
	req.Header.Set("Accept", "application/json; charset=UTF-8")
	req.Header.Set("Stream", "True")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 14096)

	for {
		n, err := resp.Body.Read(buf)
		if err != nil || n == 0 {
			break
		}
		jsonStream := string(buf[:n])
		dec := json.NewDecoder(strings.NewReader(jsonStream))
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if m.ErrorCode == 0 {
			text := strings.Split(m.Text, body.Prompt)[1]
			body.Prompt = body.Prompt + text
			data := "data: " + text + "\n\n"
			ctx.Writer.WriteString(data)
			ctx.Writer.Flush()
		}
	}

}
