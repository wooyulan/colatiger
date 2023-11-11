package response

import (
	cErr "colatiger/api/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应结构体
type Response struct {
	ErrorCode int         `json:"error_code"`
	Status    string      `json:"status"`  // 自定义错误码
	Data      interface{} `json:"data"`    // 数据
	Message   string      `json:"message"` // 信息
}

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	if gin.Mode() != gin.ReleaseMode {
		if _, ok := err.(error); ok {
			msg = err.(error).Error()
		}
	}
	FailByErr(c, cErr.InternalServer(msg))
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		"ok",
		data,
		"ok",
	})
	c.Abort()
}

func Fail(c *gin.Context, httpCode int, errorCode int, msg string) {
	c.JSON(httpCode, Response{
		errorCode,
		"fail",
		nil,
		msg,
	})
	c.Abort()
}

func FailByErr(c *gin.Context, err error) {
	v, ok := err.(*cErr.Error)
	if ok {
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		Fail(c, http.StatusBadRequest, cErr.DEFAULT_ERROR, err.Error())
	}
}
