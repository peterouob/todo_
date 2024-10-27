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
	prefix := strconv.FormatInt(id, 10)

	firstDuration := time.Until(time.Unix(tk.AtExp, 0))
	secondDuration := time.Until(time.Unix(rtk.ReExp, 0))

	firstToken := map[string]interface{}{
		"uuid":  tk.AccessUUid,
		"token": tk.AccessToken,
	}

	secondToken := map[string]interface{}{
		"uuid":  rtk.RefreshUUid,
		"token": rtk.RefreshToken,
	}

	if err := db.Rdb.HMSet(context.Background(), prefix+":access", firstToken).Err(); err != nil {
		return errors.New("error saving access token: " + err.Error())
	}
	if err := db.Rdb.Expire(context.Background(), prefix+":access", firstDuration).Err(); err != nil {
		return errors.New("error setting access token TTL: " + err.Error())
	}

	if err := db.Rdb.HMSet(context.Background(), prefix+":refresh", secondToken).Err(); err != nil {
		return errors.New("error saving refresh token: " + err.Error())
	}
	if err := db.Rdb.Expire(context.Background(), prefix+":refresh", secondDuration).Err(); err != nil {
		return errors.New("error setting refresh token TTL: " + err.Error())
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
