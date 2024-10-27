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
	"log"
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

	go func() {
		err := refreshTokenRoutine(rt.RefreshToken, id, rtokenVal)
		if err != nil {
			log.Println("error in refresh token routine :" + err.Error())
		}
	}()
	return t, rt, nil
}

// TODO:check this function can work because i change the store function which store token value to redis
func refreshTokenRoutine(refreshToken string, userId int64, refreshTokenSecret string) error {
	for {
		token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(refreshTokenSecret), nil
		})
		if err != nil {
			return errors.New("error parsing refresh token :" + err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := int64(claims["exp"].(float64))
			sleepDuration := time.Duration(exp-time.Now().Unix()-60) * time.Second
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
			if err := DeleteOldToken(userId); err != nil {
				return err
			}
			newToken, newRtoken, _ := CreateToken(userId)
			if err := SaveToken(userId, newToken, newRtoken); err != nil {
				return err
			}
		} else {
			return errors.New("error parsing refresh token :" + err.Error())
		}
	}
}
func createToken(id int64, value string) (*token2.Token, error) {
	t := &token2.Token{}
	t.AccessUUid = uuid.NewString()
	t.AtExp = time.Now().Add(time.Minute * 15).Unix()
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
	t.ReExp = time.Now().Add(time.Minute * 30).Unix()
	rclaim := jwt.MapClaims{}
	rclaim["authorized"] = true
	rclaim["refresh_uuid"] = t.RefreshUUid
	rclaim["user_id"] = id
	rclaim["exp"] = time.Now().Add(time.Hour * 30).Unix()
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

	accessToken, err := db.Rdb.HGet(context.Background(), accessTokenKey, "token").Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving access token from Redis: %v", err)
	}

	refreshToken, err := db.Rdb.HGet(context.Background(), refreshTokenKey, "token").Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving refresh token from Redis: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
		}
		if tokenString == accessToken {
			return []byte(os.Getenv("TOKENKEY")), nil

		} else if tokenString == refreshToken {
			return []byte(os.Getenv("TOKENREKEY")), nil
		}
		return nil, errors.New("not found the correct token value")
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("error parsing token or token is invalid: %v", err)
	}

	return token, nil
}
