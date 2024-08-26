package main

import (
	"github.com/WillMatthews/trump-or-markov/internal/api"
	"github.com/WillMatthews/trump-or-markov/internal/config"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, version := config.GetConfig()
	tt.LoadTrumpTweets(cfg.Dataset)
	r := gin.Default()

	tapi := api.NewTrumpAPI(&cfg.TrumpTwitter)

	r.Use(func(c *gin.Context) {
		c.Header("X-Version", version)
		c.Next()
	})

	r.GET("/v1/trump", tapi.HandleTrump)

	// basic ping/pong health check
	r.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hell yeah",
		})
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Abort()
	})

	r.Run(cfg.Server.Address())
}
