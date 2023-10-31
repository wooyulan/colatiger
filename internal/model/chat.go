package model

import "strconv"

type Chat struct {
	Id       uint64 `json:"id" gorm:"primaryKey"`
	Question string `json:"user_content" gorm:"comment:Q"`
	Answer   string `json:"answer" gorm:"comment:A"`
	File     string `json:"file" gorm:"comment:文件地址"`
	Prompt   string `json:"prompt" gorm:"comment:prompt"`
	UserId   uint64 `json:"user_id" gorm:"comment:用户id"`
	Timestamps
	SoftDeletes
}

func (c *Chat) TableName() string {
	return "t_chat_his"
}

func (c *Chat) GetUid() string {
	return strconv.Itoa(int(c.Id))
}
