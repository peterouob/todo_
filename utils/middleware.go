package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func Cors(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	fmt.Println(c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	c.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func ValidApiKey(c *gin.Context) {
	if c.Request.Header.Get("API-KEY") != os.Getenv("APIKEY") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -1,
			"msg":  "needed api key",
		})
		c.Abort()
	}
	c.Next()
}

func AuthByJWT() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": "-1",
				"msg:": "not have auth header",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Format of Authorization is wrong",
			})
			c.Abort()
			return
		}
		_, err := VerifyToken(c, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "Verify of Authorization is wrong :" + err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
