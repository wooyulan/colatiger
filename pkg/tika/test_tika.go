package main

import (
	"colatiger/pkg/helper/readhtml"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-tika/tika"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "fp", "/Users/eric/Downloads/测试正常文字.pdf", "pdf file path.")
	flag.Parse()
	if filePath == "" {
		panic("file path must be provided")
	}
	content, err := readPdf(filePath) // Read local pdf file
	if err != nil {
		panic(err)
	}
	// fmt.Println(content)

	//将pdf的所有内容写入html文件)
	err = os.WriteFile("./out.html", []byte(content), 0666)
	if err != nil {
		log.Fatal(err)
	}

	//先将html中的<title>标签去掉,因为此标签中含有特殊字符,会导致xml语法出错
	delerr := deleteTitle("out.html")
	if delerr != nil {
		log.Fatal(delerr)
	}

	err = readhtml.ReadHtml("out.html")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func readPdf(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	client := tika.NewClient(nil, "http://127.0.0.1:9998")
	return client.Parse(context.TODO(), f)
}

// 删除html中的title标签
func deleteTitle(filename string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("sed -i '65d' %s", filename))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
