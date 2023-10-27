package img_base

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

func GetUrlImgBase64(path string) (baseImg string, err error) {
	//获取网络图片
	client := &http.Client{
		Timeout: time.Second * 5, //超时时间
	}
	var bodyImg io.Reader
	request, err := http.NewRequest("GET", path, bodyImg)
	if err != nil {
		err = errors.New("获取网络图片失败")
		return
	}
	respImg, _ := client.Do(request)
	defer respImg.Body.Close()
	imgByte, _ := io.ReadAll(respImg.Body)
	baseImg = base64.StdEncoding.EncodeToString(imgByte)

	return baseImg, nil
}

func GetUrlImgBase64Prefix(path string) (baseImg string, err error) {
	//获取网络图片
	client := &http.Client{
		Timeout: time.Second * 5, //超时时间
	}
	var bodyImg io.Reader
	request, err := http.NewRequest("GET", path, bodyImg)
	if err != nil {
		err = errors.New("获取网络图片失败")
		return
	}
	respImg, _ := client.Do(request)
	defer respImg.Body.Close()
	imgByte, _ := io.ReadAll(respImg.Body)

	// 判断文件类型，生成一个前缀，拼接base64后可以直接粘贴到浏览器打开，不需要可以不用下面代码
	//取图片类型
	mimeType := http.DetectContentType(imgByte)
	switch mimeType {
	case "image/jpeg":
		baseImg = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imgByte)
	case "image/png":
		baseImg = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgByte)
	}
	return
}
