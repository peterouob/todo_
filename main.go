package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/peterouob/todo_/mysql"
	"github.com/peterouob/todo_/router"
	"log"
	"net/http"
)

//
//func init() {
//	if err := godotenv.Load(); err != nil {
//		log.Panicf("error to load env file ... :%s", err.Error())
//	}
//}

func main() {
	go func() {
		mysql.InitMysql()
	}()

	r := gin.Default()
	r.Use(Cors)
	router.InitRouter(r)
	if err := r.Run(":8084"); err != nil {
		log.Panicf("errors:%s", err.Error())
	}
}

func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	fmt.Println(c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	c.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
