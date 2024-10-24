package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) TableName() string {
	return "user"
}
