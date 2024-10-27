package utils

import (
	"github.com/gin-gonic/gin"
	"os"
)

func SetCookie(c *gin.Context, name, value string) {
	c.SetCookie(name, value, 365*3600, "/", os.Getenv("SERVER"), false, true)
}

//func SetCookie(c *gin.Context, name, value string) {
//	c.SetCookie(name, value, 365*3600, "/", "localhost", false, true)
//}

func RemoveCookie(c *gin.Context, key string) {
	c.SetCookie(key, "", -1, "", os.Getenv("SERVER"), false, true)
}
