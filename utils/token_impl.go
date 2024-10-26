package utils

import (
	"context"
	"errors"
	"github.com/peterouob/todo_/db"
	"github.com/peterouob/todo_/model"
	"strconv"
	"time"
)

func SaveToken(id int64, tk *model.Token, rtk *model.RefreshToken) error {
	prefix := id

	firstTime := time.Unix(tk.AtExp, 0)
	secondTime := time.Unix(rtk.ReExp, 0)
	now := time.Now()
	firstToken := map[string]interface{}{
		"uuid":  tk.AccessUUid,
		"token": tk.AccessToken,
		"time":  firstTime.Sub(now),
	}

	secondToken := map[string]interface{}{
		"uuid":  rtk.RefreshUUid,
		"token": rtk.RefreshToken,
		"time":  secondTime.Sub(now),
	}

	if err := db.Rdb.HMSet(context.Background(), strconv.FormatInt(prefix, 10), firstToken, secondToken).Err(); err != nil {
		return errors.New("error in save token :" + err.Error())
	}
	return nil
}

func DeleteOldToken(key int64) error {
	if err := db.Rdb.HDel(context.Background(), strconv.FormatInt(key, 10)).Err(); err != nil {
		return errors.New("error in delete old token")
	}
	return nil
}

func GetTokenMapById(id int64) []interface{} {
	result := db.Rdb.HMGet(context.Background(), strconv.FormatInt(id, 10))
	return result.Val()
}
