package router

import "github.com/gin-gonic/gin"

func InitRouter(r *gin.Engine) {
	r.GET("/")
	r.POST("/login")
	r.POST("/register")
	r.POST("/update")
	r.POST("/delete")
}
