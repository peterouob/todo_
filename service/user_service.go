package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/peterouob/todo_/model"
	"github.com/peterouob/todo_/mysql"
	"gorm.io/gorm"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	user := model.User{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "bind user data error:" + err.Error(),
		})
	}

	if err := mysql.DB.Where("username=?", user.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		mysql.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "welcome",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "have same user!",
		})
		return
	}
}

func LoginUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "bind user data error:" + err.Error(),
		})
	}

	if err := mysql.DB.Where("username=? AND password=?", user.Username, user.Password).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "error in login",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success login",
	})
}
