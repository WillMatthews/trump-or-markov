package main

import (
	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

const (
	PORT = "1234"
)

func initialise(cfg *config.Config) {
	trumptweets.LoadTrumpTweets(cfg.Dataset)
}

func main() {
	cfg := config.GetConfig()
	initialise(cfg)

	r := gin.Default()

	r.GET("/realDTTweet", func(c *gin.Context) {
		tweet := trumptweets.RandomSample()
		c.JSON(200, tweet)
	})

	r.GET("/fakeDTTweet", func(c *gin.Context) {
		tweet := trumptweets.RandomFakeSample(2)
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
