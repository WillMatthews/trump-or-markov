package api

import (
	"strconv"

	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

func Trump(c *gin.Context) {
	ord := 2
	if ordQry, ok := c.GetQuery("ord"); ok {
		if parsed, err := strconv.Atoi(ordQry); err == nil {
			ord = parsed
		}
	}

	numTweets := 1
	if num, ok := c.GetQuery("n"); ok {
		if parsed, err := strconv.Atoi(num); err == nil {
			numTweets = parsed
		}
	}

	if makeFake, ok := c.GetQuery("fake"); ok {
		if makeFake == "true" {
			fake(c, ord, numTweets)
			return
		}
	}

	real(c, numTweets)
}

func fake(c *gin.Context, ord int, num int) {
	createTweets(c, num, func() (*tt.Tweet, error) {
		return tt.RandomFakeSample(ord)
	})
}

func real(c *gin.Context, num int) {
	createTweets(c, num, tt.RandomSample)
}

func createTweets(c *gin.Context, num int, f func() (*tt.Tweet, error)) {
	var tweets []*tt.Tweet
	for i := 0; i < num; i++ {
		tweet, err := f()
		check(c, err)
		tweets = append(tweets, tweet)
	}
	c.JSON(200, tweets)
}

func check(c *gin.Context, err error) {
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
}
