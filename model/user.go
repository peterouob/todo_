package model

import (
	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primarykey"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
