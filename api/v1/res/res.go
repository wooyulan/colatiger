package res

import (
	"colatiger/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// 响应结构体
type Response struct {
	Status  string      `json:"status"`  // 自定义错误码
	Data    interface{} `json:"data"`    // 数据
	Message string      `json:"message"` // 信息
}

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		"ok",
		data,
		"ok",
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode string, msg string) {
	c.JSON(http.StatusOK, Response{
		errorCode,
		nil,
		msg,
	})
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error v1.CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, v1.Errors.ValidateError.ErrorCode, msg)
}

// BusinessFail 业务逻辑失败
func BusinessFail(c *gin.Context, msg string) {
	Fail(c, v1.Errors.BusinessError.ErrorCode, msg)
}

func TokenFail(c *gin.Context) {
	FailByError(c, v1.Errors.TokenError)
}

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	// 非生产环境显示具体错误信息
	if os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
		if _, ok := err.(error); ok {
			msg = err.(error).Error()
		}
	}
	c.JSON(http.StatusInternalServerError, Response{
		"fail",
		nil,
		msg,
	})
	c.Abort()
}