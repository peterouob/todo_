package router

import (
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/service"
	"github.com/peterouob/todo_/utils"
)

func InitRouter(r *gin.Engine) {
	r.Use(utils.Cors, utils.ValidApiKey)
	r.GET("/")
	r.POST("/login", service.LoginUser)
	r.POST("/register", service.RegisterUser)
	todo := r.Group("/todo")
	todo.Use(utils.AuthByJWT())
	{
		todo.GET("/", service.GetAllTodo)
		todo.GET("/:id", service.GetTodoByID)
		todo.GET("/filter", service.FilterByMonthAndYear)
		todo.POST("/done", service.GetTodoFilterDone)
		todo.POST("/create", service.CreateTodo)
		todo.PUT("/:id", service.UpdateTodo)
		todo.DELETE("/:id", service.DeleteTodo)
		todo.GET("/done/:id", service.Done)
	}
}
