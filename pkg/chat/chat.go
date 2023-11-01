package chat

import (
	"bufio"
	"bytes"
	v1 "colatiger/api/v1/req"
	"colatiger/internal/model"
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

func SendMsg(ctx *gin.Context, llavaReq *LLaVaChatReq) string {
	jsonBytes, err := json.Marshal(llavaReq)
	if err != nil {
	}
	req, err := http.NewRequest("POST", "http://82.156.138.158:10000/worker_generate_stream", bytes.NewBuffer(jsonBytes))
	if err != nil {
		println(err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	defer ctx.Writer.Flush()

	var assistant bytes.Buffer
	reader := bufio.NewReader(resp.Body)
	for {
		// 读取到分隔符0
		line, err := reader.ReadString(byte(0))
		if err != nil || line == "" {
			break
		}
		// 去掉分隔符0
		line = strings.Trim(line, "\x00")
		dec := json.NewDecoder(strings.NewReader(line))
		var m MessageResponse
		if err := dec.Decode(&m); err == io.EOF || err != nil {
			break
		}
		if m.ErrorCode == 0 {
			text := strings.Split(m.Text, llavaReq.Prompt)[1]
			llavaReq.Prompt = llavaReq.Prompt + text
			assistant.WriteString(text)
			textTemp, _ := json.Marshal(text)
			var data = "data: " + string(textTemp) + "\n\n"
			fmt.Fprintf(ctx.Writer, data)
			ctx.Writer.Flush()
		}
		if resp.Body == nil {
			break
		}

	}
	return assistant.String()
}

func BuildLLaVaModelBody(ctx *gin.Context, chatReq v1.ChatReq, his *[]model.Chat) string {

	body := &LLaVaChatReq{
		Model:        "llava-v1.5-13b",
		Prompt:       "",
		MaxNewTokens: 512,
		Temperature:  0.7,
		Stop:         "</s>",
		Images:       []string{},
	}
	prompt := BuildHisMessage(his, body)
	body.Prompt = prompt

	if chatReq.Images != nil && len(chatReq.Images) > 0 {
		// base64 图片
		for _, imgUrl := range chatReq.Images {
			base64, _ := img.GetUrlImgBase64(imgUrl)
			body.Images = append(body.Images, base64)
		}
	}

	if len(body.Images) > 0 {
		body.Prompt = fmt.Sprintf(prompt, "<image>\\n"+chatReq.Message)
	} else {
		body.Prompt = fmt.Sprintf(prompt, chatReq.Message)
	}

	assistantRes := SendMsg(ctx, body)
	return assistantRes
}

func BuildHisMessage(his *[]model.Chat, req *LLaVaChatReq) string {
	prompt := "A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions. "

	//prompt := "\"\"\"\n- 你会将图片通过OCR后的文本信息整合总结，请一步一步思考，你会挖掘不同单词和信息之间的联系，\n你会用各种信息分析方法（如：统计、聚类...等）完成信息整理任务，翻译成中文回复。\n- 输出格式：\n  \"\n  # 图片的信息\n  # 图片描述\n  {填充信息：通过一句话描述整体内容，不超过30字}\n  {填充信息：分点显示，整合信息后总结，最多不超过8点，每条信息不超过20字，保留关键值，如人名、地名...}\n  # 分类\n  {填充信息：为该信息3-5个分类标签，例如：#身份证、#银行卡、#驾驶证、#护照}\n  \"\n\"\"\""

	hisPrompt := ""
	template := "USER: %s ASSISTANT:%s </s>"

	for _, msg := range *his {
		if ok := msg.File != ""; ok {
			base64, _ := img.GetUrlImgBase64(msg.File)
			req.Images = append(req.Images, base64)
			//hisPrompt += fmt.Sprintf(template, "<image>\\n"+msg.Question, msg.Answer)
		}
		hisPrompt += fmt.Sprintf(template, msg.Question, msg.Answer)
	}
	return prompt + hisPrompt + "USER: %s ASSISTANT: "
}
