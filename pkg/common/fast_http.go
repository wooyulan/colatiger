package common

import (
	"crypto/tls"
	"encoding/json"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/valyala/fasthttp"
	"net/url"
	"time"
)

var FastClient = fastClient()

func fastClient() fasthttp.Client {
	return fasthttp.Client{
		Name:                     "ColaTiger_Client",
		NoDefaultUserAgentHeader: true,
		TLSConfig:                &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:          2000,
		MaxIdleConnDuration:      5 * time.Second,
		MaxConnDuration:          5 * time.Second,
		ReadTimeout:              5 * time.Second,
		WriteTimeout:             5 * time.Second,
		MaxConnWaitTimeout:       5 * time.Second,
	}
}

func SendPost(headers map[string]string, reqUri string, data url.Values) (resBody *jsonvalue.V, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	// 设置header
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// 设置url
	req.Header.SetRequestURI(reqUri)
	// 设置body
	req.SetBodyString(data.Encode())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := FastClient.Do(req, resp); err != nil {
		return nil, err
	}

	j, err := jsonvalue.Unmarshal(resp.Body())
	if err != nil {
		return nil, err
	}
	if code, err := j.GetString("code"); err != nil || code != "0" {
		return nil, err
	}
	return j, nil
}

func SendSimplePost(reqUri string, data map[string]interface{}) (resBody *jsonvalue.V, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")

	body, _ := json.Marshal(data)
	req.SetBody(body)
	// 设置url
	req.Header.SetRequestURI(reqUri)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	if err := FastClient.Do(req, resp); err != nil {
		return nil, err
	}
	j, err := jsonvalue.Unmarshal(resp.Body())
	if err != nil {
		return nil, err
	}
	return j, nil
}
