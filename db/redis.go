package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
)

var Rdb *redis.Client

func InitRedis() {
	db, _ := strconv.Atoi(os.Getenv("RDB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("RDBADDR"), os.Getenv("RDBPORT")),
		Password: os.Getenv("RDBPASS"),
		DB:       db,
	})
	Rdb = rdb
	log.Println("connect redis success ...")
}
