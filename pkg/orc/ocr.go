package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	//	img.DownPic("https://c-ssl.duitang.com/uploads/item/201506/18/20150618012527_urPA8.jpeg", "/Users/eric/Desktop/test")

	client.SetImage("/Users/eric/Desktop/2.jpeg")

	client.SetLanguage("chi_sim")
	text, _ := client.Text()
	fmt.Println(text)
}
