package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/service"
)

func InitRouter(r *gin.Engine) {
	r.GET("/")
	r.POST("/login", service.LoginUser)
	r.POST("/register", service.RegisterUser)
}
