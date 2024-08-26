package api

import (
	"errors"
	"strconv"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	tt "github.com/WillMatthews/trump-or-markov/internal/trumptweets"
	"github.com/gin-gonic/gin"
)

var (
	ErrOrderTooHigh = errors.New("requested markov order is too high")
	ErrNumTweets    = errors.New("requested number of tweets is invalid or too high")
)

func Trump(c *gin.Context, config *config.TrumpTwitter) {
	ord, err := parseOrd(c, config.Markov.MaxOrder)
	if err != nil {
		sendError(c, err)
		return
	}

	numTweets, err := parseNumTweets(c, config.MaxTweets)
	if err != nil {
		sendError(c, err)
		return
	}

	if makeFake, ok := c.GetQuery("fake"); ok {
		if makeFake == "true" {
			fake(c, ord, numTweets, config)
			return
		}
	}

	real(c, numTweets, &config.Markov)
}

func parseOrd(c *gin.Context, maxOrder int) (int, error) {
	ord := 2
	if ordQry, ok := c.GetQuery("ord"); ok {
		if parsed, err := strconv.Atoi(ordQry); err == nil {
			ord = parsed
		}
	}

	if ord > maxOrder {
		return 0, ErrOrderTooHigh
	}
	return ord, nil
}

func parseNumTweets(c *gin.Context, maxTweets int) (int, error) {
	numTweets := 1
	if num, ok := c.GetQuery("n"); ok {
		if parsed, err := strconv.Atoi(num); err == nil {
			numTweets = parsed
		}
	}

	if numTweets < 1 || numTweets > maxTweets {
		return 0, ErrNumTweets
	}
	return numTweets, nil
}

func fake(c *gin.Context,
	markovOrder int,
	numTweets int,
	config *config.TrumpTwitter,
) {
	createTweets(c, numTweets, func() (*tt.Tweet, error) {
		return tt.RandomFakeSample(markovOrder, config)
	})
}

func real(c *gin.Context,
	num int,
	cfg *config.Markov,
) {
	gen := func() (*tt.Tweet, error) {
		return tt.RandomRealSample(cfg)
	}

	createTweets(c, num, gen)
}

func createTweets(c *gin.Context, num int, tweetGen func() (*tt.Tweet, error)) {
	var tweets []*tt.Tweet
	for i := 0; i < num; i++ {
		tweet, err := tweetGen()
		check(c, err)
		tweets = append(tweets, tweet)
	}
	c.JSON(200, tweets)
}

func sendError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}

func check(c *gin.Context, err error) {
	if err != nil {
		sendError(c, err)
	}
}
