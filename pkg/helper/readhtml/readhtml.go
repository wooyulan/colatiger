package readhtml

import (
	"encoding/xml"
	"os"
)

// 定义多个结构体，用来接收反序列化数据
type (
	html struct {
		XMLNAME xml.Name `xml:"html"`
		Body    htmlBody `xml:"body"`
	}
	htmlBody struct {
		XMLNAME xml.Name `xml:"body"`
		Div     htmlDiv  `xml:"div"`
	}
	htmlDiv struct {
		P []string `xml:"p"`
	}
)

// 读取html文件，并反序列化到结构体中
func ReadHtml(filename string) error {
	rf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	html := html{}
	err = xml.Unmarshal(rf, &html)
	if err != nil {
		return err
	}
	var b []byte
	for _, v := range html.Body.Div.P {
		b = append(b, []byte(v)...)
	}
	err = os.WriteFile("res.doc", b, 0644)
	if err != nil {
		return err
	}
	return nil
}
