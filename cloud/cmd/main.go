package main

import (
	"github.com/WillMatthews/trump-or-markov/internal/api"
	"github.com/WillMatthews/trump-or-markov/internal/config"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

const (
	PORT = "1776" // why not, it's very apt. Only really used by FEMA so should be fine.
)

func initialise(cfg *config.Config) {
	tt.LoadTrumpTweets(cfg.Dataset)
}

func main() {
	cfg, version := config.GetConfig()
	initialise(cfg)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("X-Version", version)
		c.Next()
	})

	r.GET("/v1/trump", func(c *gin.Context) {
		api.Trump(c, &cfg.TrumpTwitter)
	})

	// basic ping/pong health check
	r.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hell yeah",
		})
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Abort()
	})

	r.Run(":" + PORT)
}
