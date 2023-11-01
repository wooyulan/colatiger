package strategy

import (
	"colatiger/pkg/common"
	"colatiger/pkg/helper/img"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"strconv"
	"time"
)

const (
	appid       = ""
	apikey      = ""
	reqUrl      = "https://webapi.xfyun.cn/v1/service/v1/ocr/idcard"
	contentType = "application/x-www-form-urlencoded"
)

type IDCard struct {
	Address        string `json:"address"`
	Birthday       string `json:"birthday"`
	IdNumber       string `json:"id_number"`
	Name           string `json:"name"`
	People         string `json:"people"`
	Sex            string `json:"sex"`
	Type           string `json:"type"`
	IssueAuthority string `json:"issue_authority"`
	Validity       string `json:"validity"`
}

func (*IDCard) Ocr(ctx *OcrContext) (interface{}, error) {
	header, data, _ := builderIdCardReq(ctx.ImgUrl)

	res, err := common.SendPost(header, reqUrl, data)

	if err != nil {
		return nil, errors.New("调用上游服务异常")
	}
	idcard, _ := res.GetString("data", "id_number")
	name, _ := res.GetString("data", "name")
	birthday, _ := res.GetString("data", "birthday")
	people, _ := res.GetString("data", "people")
	sex, _ := res.GetString("data", "sex")
	types, _ := res.GetString("data", "type")
	address, _ := res.GetString("data", "address")
	issue_authority, _ := res.GetString("data", "issue_authority")
	validity, _ := res.GetString("data", "validity")

	var idCard = &IDCard{
		IdNumber:       idcard,
		Name:           name,
		Address:        address,
		Birthday:       birthday,
		People:         people,
		Sex:            sex,
		Type:           types,
		IssueAuthority: issue_authority,
		Validity:       validity,
	}
	// 解析返回值
	return idCard, nil
}

func builderIdCardReq(imgUrl string) (map[string]string, url.Values, error) {
	// 组装请求参数
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	param := make(map[string]string)
	// 引擎类型930820
	param["engine_type"] = "idcard"
	// 是否返回头像图片
	param["head_portrait"] = "0"
	tmp, _ := json.Marshal(param)
	base64_param := base64.StdEncoding.EncodeToString(tmp)
	w := md5.New()
	io.WriteString(w, apikey+curtime+base64_param)
	checksum := fmt.Sprintf("%x", w.Sum(nil))

	f_base64, err := img.GetUrlImgBase64(imgUrl)

	if err != nil {
		return nil, nil, err
	}

	data := url.Values{}
	data.Add("image", f_base64)

	// 组装http请求头
	header := make(map[string]string)
	header["X-Appid"] = appid
	header["X-CurTime"] = curtime
	header["X-Param"] = base64_param
	header["X-CheckSum"] = checksum
	header["Content-Type"] = contentType

	return header, data, nil
}
