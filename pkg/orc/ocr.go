package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("/Users/eric/Desktop/test.jpeg")
	client.SetLanguage("chi_sim")
	text, _ := client.Text()
	fmt.Println(text)
}
