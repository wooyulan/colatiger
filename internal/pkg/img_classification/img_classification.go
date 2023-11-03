package img_classification

import (
	v1 "colatiger/api/v1/req"
	"colatiger/pkg/common"
)

const imgUrl = "http://localhost:5050/open/v1/img"

type ImgClassification struct {
	ImgUrl  string `json:"img_url"`
	ImgType string `json:"imgType"`
}

// 发起远程调用
func ImgClassificationPost(req v1.ImgClassificationReq) (*ImgClassification, error) {

	reqBody := map[string]interface{}{
		"filePath": req.ImgUrl,
	}

	res, err := common.SendSimplePost(imgUrl, reqBody)
	if err != nil {
		return nil, err
	}

	imgType, _ := res.GetString("type")

	return &ImgClassification{
		ImgUrl:  req.ImgUrl,
		ImgType: imgType,
	}, nil

}
