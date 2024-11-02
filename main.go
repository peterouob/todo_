package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/peterouob/todo_/db"
	"github.com/peterouob/todo_/router"
	"log"
	"os"
)

func main() {
	go func() {
		db.InitMysql()
		db.InitMongo()
		db.InitRedis()
	}()

	r := gin.Default()
	router.InitRouter(r)
	if err := r.Run(os.Getenv("PORT")); err != nil {
		log.Panicf("errors:%s", err.Error())
	}
}
