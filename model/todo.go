package model

import "time"

type Todo struct {
	Title    string    `json:"title" bson:"title"`
	Content  string    `json:"content" bson:"content"`
	DeadTime time.Time `json:"dead_time" bson:"deadTime"`
	Done     bool      `json:"done" bson:"done"`
}
