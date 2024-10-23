package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "GET METHOD",
		})
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "POST METHOD",
		})
	})
	r.Run(":8084")
}
