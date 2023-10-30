package main

/**
 *身份证识别 WebAPI 接口调用示例 接口文档（必看）：https://doc.xfyun.cn/rest_api/%E8%BA%AB%E4%BB%BD%E8%AF%81%E8%AF%86%E5%88%AB.html
 *webapi OCR服务参考帖子（必看）：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=39111&highlight=OCR
 *(Very Important)创建完webapi应用添加身份证识别之后一定要设置ip白名单，找到控制台--我的应用--设置ip白名单，如何设置参考：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=41891
 *图片属性：仅支持jpg格式，推荐 jpg 文件设置为：尺寸 1024×768，图像质量 75 以上，位深度 24。base64位编码之后大小不超过4M
 *错误码链接：https://www.xfyun.cn/document/error-code (code返回错误码时必看)
 *OCR错误码400开头请在接口文档底部查看
 */
import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func ocr_idcard() {
	// 应用APPID(必须为webapi类型应用，并开通身份证识别服务，参考帖子如何创建一个webapi应用：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=36481)
	appid := "04ada1cb"
	// 接口密钥(webapi类型应用开通身份证识别服务后,控制台--我的应用---身份证识别---相应服务的apikey)
	apikey := "31952d61759f8c907676672cfbd15800"
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
	// 上传图片地址
	f, _ := os.ReadFile("/Users/eric/Desktop/2.jpg")
	f_base64 := base64.StdEncoding.EncodeToString(f)
	data := url.Values{}
	data.Add("image", f_base64)
	body := data.Encode()
	client := &http.Client{}
	// 身份证识别webapi接口地址
	req, _ := http.NewRequest("POST", "https://webapi.xfyun.cn/v1/service/v1/ocr/idcard", strings.NewReader(body))
	// 组装http请求头
	req.Header.Set("X-Appid", appid)
	req.Header.Set("X-CurTime", curtime)
	req.Header.Set("X-Param", base64_param)
	req.Header.Set("X-CheckSum", checksum)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(req)
	defer res.Body.Close()
	res_body, _ := io.ReadAll(res.Body)
	fmt.Print(string(res_body))
}

func main() {
	ocr_idcard()
}
