package v1

type Register struct {
	Name            string `form:"name" json:"name" binding:"required"`
	Mobile          string `form:"mobile" json:"mobile"`
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirmPassword" json:"confirmPassword" binding:"required"`
}

// 自定义错误信息
func (register Register) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "用户名称不能为空",
		"mobile.required":   "手机号码不能为空",
		"email.required":    "邮箱不能为空",
		"password.required": "用户密码不能为空",
	}
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (login Login) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"mobile.required":   "手机号码不能为空",
		"mobile.mobile":     "手机号码格式不正确",
		"email.required":    "邮箱不能为空",
		"password.required": "用户密码不能为空",
	}
}
