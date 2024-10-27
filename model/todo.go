package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"
)

type Todo struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	DeadTime time.Time          `json:"dead_time" bson:"deadTime"`
	Done     bool               `default:"false" json:"done" bson:"done"`
	UserID   int64              `json:"user_id" bson:"userID"`
}

type TodoGroup struct {
	GroupID struct {
		Year  int `bson:"year"`
		Month int `bson:"month"`
	} `bson:"_id"`
	Todos []Todo `json:"todos" bson:"todos"`
}

type FindTodoRequest struct {
	Done bool `default:"false" json:"done" bson:"done,omitempty"`
}

type UpdateTodoRequest struct {
	Title    string    `json:"title" bson:"title"`
	Content  string    `json:"content" bson:"content,omitempty"`
	DeadTime time.Time `json:"dead_time" bson:"deadTime,omitempty"`
}

func NewTodo() Todo {
	todo := Todo{}
	todo.Id = primitive.NewObjectID()
	setStructDefaultTag(&todo)
	return todo
}

func setStructDefaultTag(s interface{}) {
	val := reflect.ValueOf(s).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if defaultTag, ok := field.Tag.Lookup("default"); ok && defaultTag == "false" {
			if val.Field(i).Kind() == reflect.Bool && !val.Field(i).Bool() {
				val.Field(i).SetBool(false)
			}
		}
	}
}
