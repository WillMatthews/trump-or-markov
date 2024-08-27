package main

import (
	"time"

	"github.com/WillMatthews/trump-or-markov/internal/api"
	"github.com/WillMatthews/trump-or-markov/internal/config"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	cfg, _ := config.GetConfig()
	tt.LoadTrumpTweets(cfg.Dataset)
	r := gin.Default()

	tapi := api.NewTrumpAPI(&cfg.TrumpTwitter)

	r.Use(cors.New(corsConfig))

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
