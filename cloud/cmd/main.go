package main

import (
	"github.com/WillMatthews/trump-or-markov/internal/api"
	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/WillMatthews/trump-or-markov/internal/logging"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

const (
	PORT    = "1776" // why not, it's very apt. Only really used by FEMA so should be fine.
	logJSON = false
)

var (
	conf, version = config.GetConfig()
)

func init() {
	logging.SetupLogger(logJSON)
	tt.LoadTrumpTweets(conf.Dataset)
	gin.SetMode(gin.DebugMode)
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("X-Version", version)
		c.Next()
	})

	r.GET("/v1/trump", func(c *gin.Context) {
		api.Trump(c)
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
