package service

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/model"
	"github.com/peterouob/todo_/utils"
	"net/http"
	"strconv"
)

func RegisterUser(c *gin.Context) {
	user := model.User{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  errors.New("username or password cannot be empty"),
		})
		return
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	user.ID = node.Generate().Int64()
	if err := registerUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "welcome",
	})
}

func LoginUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  errors.New("username or password cannot be empty"),
		})
		return
	}
	uid, err := loginUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	user.ID = uid

	utils.SetCookie(c, "id", strconv.FormatInt(uid, 10))
	tk, rtk, err := utils.CreateToken(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":              -1,
			"create token err:": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          0,
		"msg":           "login success",
		"data":          user,
		"token":         tk,
		"refresh_token": rtk,
	})
}
