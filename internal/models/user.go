package models

type User struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;comment:用户名称"`
	Mobile   string `json:"mobile" gorm:"not null;index;comment:用户手机号"`
	Password string `json:"password" gorm:"not null;default:'';comment:用户密码"`
	Timestamps
	SoftDeletes
}

func (u *User) TableName() string {
	return "t_user"
}
