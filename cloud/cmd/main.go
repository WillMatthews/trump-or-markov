package main

import (
	"strconv"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

const (
	PORT = "1234"
)

func initialise(cfg *config.Config) {
	tt.LoadTrumpTweets(cfg.Dataset)
}

func main() {
	cfg := config.GetConfig()
	initialise(cfg)
	r := gin.Default()

	r.GET("/realDTTweet", func(c *gin.Context) {
		tweet, err := tt.RandomSample()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, tweet)
	})

	r.GET("/fakeDTTweet", func(c *gin.Context) {
		ord := 2
		if ordQry, ok := c.GetQuery("order"); ok {
			if parsed, err := strconv.Atoi(ordQry); err == nil {
				ord = parsed
			}
		}

		tweet, err := tt.RandomFakeSample(ord)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

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
