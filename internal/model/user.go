package model

import "strconv"

type User struct {
	Id       uint64 `json:"id" gorm:"primaryKey"`
	Mobile   string `json:"mobile" gorm:"comment:用户手机号"`
	Nickname string `json:"nickname" gorm:"comment:用户昵称"`
	Avatar   string `json:"avatar" gorm:"comment:用户头像"`
	Email    string `json:"email" gorm:"not null;index;comment:邮箱"`
	Password string `json:"password" gorm:"not null;default:'';comment:用户密码"`
	Timestamps
	SoftDeletes
}

func (u *User) TableName() string {
	return "t_user"
}

func (u *User) GetUid() string {
	return strconv.Itoa(int(u.Id))
}
