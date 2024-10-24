package model

import (
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	gorm.Model
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Done     bool      `json:"done"`
	DeadLine time.Time `json:"date_time"`
}

func (t *Todo) TableName() string {
	return "todo"
}
