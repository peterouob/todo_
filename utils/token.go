package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/peterouob/todo_/db"
	token2 "github.com/peterouob/todo_/model"
	"os"
	"strconv"
	"time"
)

var err error
var (
	auuid chan interface{}
	ruuid chan interface{}
)

func init() {
	auuid = make(chan interface{}, 1024)
	ruuid = make(chan interface{}, 1024)
}

func CreateToken(id int64) (*token2.Token, *token2.RefreshToken, error) {
	tokenVal := os.Getenv("TOKENKEY")
	rtokenVal := os.Getenv("TOKENREKEY")
	t, err := createToken(id, tokenVal)
	if err != nil {
		return nil, nil, err
	}
	rt, err := createRefreshToken(id, rtokenVal)
	if err != nil {
		return nil, nil, err
	}
	if err := SaveToken(id, t, rt); err != nil {
		return nil, nil, err
	}
	return t, rt, nil
}
func createToken(id int64, value string) (*token2.Token, error) {
	t := &token2.Token{}
	t.AccessUUid = uuid.NewString()
	t.AtExp = time.Now().Add(time.Hour * 2).Unix()
	claim := jwt.MapClaims{}
	claim["authorized"] = true
	claim["access_uuid"] = t.AccessUUid
	claim["user_id"] = id
	claim["exp"] = t.AtExp
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t.AccessToken, err = tk.SignedString([]byte(value))
	if err != nil {
		fmt.Println("sign token error: ", err)
		return nil, err
	}
	auuid <- claim["access_uuid"]
	return t, nil
}
func createRefreshToken(id int64, value string) (*token2.RefreshToken, error) {
	t := &token2.RefreshToken{}
	t.RefreshUUid = uuid.NewString()
	t.ReExp = time.Now().Add(time.Hour * 128).Unix()
	rclaim := jwt.MapClaims{}
	rclaim["authorized"] = true
	rclaim["refresh_uuid"] = t.RefreshUUid
	rclaim["user_id"] = id
	rclaim["exp"] = time.Now().Add(time.Hour * 128).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rclaim)
	t.RefreshToken, err = rt.SignedString([]byte(value))
	if err != nil {
		return nil, errors.New("sign refresh token error:" + err.Error())
	}
	ruuid <- rclaim["refresh_uuid"]
	return t, nil
}

// VerifyToken 檢查token簽名方法
func VerifyToken(c *gin.Context, tokenString string) (*jwt.Token, error) {
	uidStr, err := c.Cookie("id")
	if err != nil {
		return nil, errors.New("get cookie value error: " + err.Error())
	}
	id, err := strconv.ParseInt(uidStr, 0, 64)
	if err != nil {
		return nil, errors.New("parse string to int64 error: " + err.Error())
	}
	accessTokenKey := fmt.Sprintf("%d:access", id)
	refreshTokenKey := fmt.Sprintf("%d:refresh", id)
	accessToken, _ := db.Rdb.HGet(context.Background(), accessTokenKey, "token").Result()
	refreshToken, _ := db.Rdb.HGet(context.Background(), refreshTokenKey, "token").Result()
	token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
		}
		if tokenString == accessToken {
			return []byte(os.Getenv("TOKENKEY")), nil

		}
		if tokenString == refreshToken {
			return []byte(os.Getenv("TOKENREKEY")), nil
		}
		return nil, errors.New("not found the correct token value")
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token or token is invalid: %v", err)
	}
	return token, nil
}
