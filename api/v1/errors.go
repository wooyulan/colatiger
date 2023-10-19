package v1

type CustomError struct {
	ErrorCode string
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError CustomError
	ValidateError CustomError
	TokenError    CustomError
	DataError     CustomError
}

var Errors = CustomErrors{
	BusinessError: CustomError{"40000", "业务错误"},
	ValidateError: CustomError{"42200", "请求参数错误"},
	TokenError:    CustomError{"40100", "登录授权失效"},
	DataError:     CustomError{"40001", "数据异常或者不存在"},
}
