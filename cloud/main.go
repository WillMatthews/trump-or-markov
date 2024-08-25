package main

import (
	"github.com/gin-gonic/gin"
)

const (
	PORT = "1234"
)

func main() {
	r := gin.Default()

	r.GET("/realDTTweet", func(c *gin.Context) {
		tweet := real.RandomSample()
		c.JSON(200, tweet)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Abort()
	})

	r.Run(":" + PORT)
}
