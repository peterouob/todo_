package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64  `gorm:"primarykey"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
