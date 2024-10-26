package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"testing"
)

func TestConnectRedis(t *testing.T) {
	db, _ := strconv.Atoi(os.Getenv("RDB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("RDBADDR"), os.Getenv("RDBPORT")),
		Password: os.Getenv("RDBPASS"),
		DB:       db,
	})
	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		t.Logf("error:%s", err.Error())
	}
}
