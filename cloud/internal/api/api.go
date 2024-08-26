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

type TrumpAPI struct {
	config *config.TrumpTwitter
}

func NewTrumpAPI(config *config.TrumpTwitter) *TrumpAPI {
	return &TrumpAPI{
		config: config,
	}
}

func (api *TrumpAPI) HandleTrump(c *gin.Context) {
	ord, err := api.parseOrd(c)
	if err != nil {
		api.sendError(c, err)
		return
	}

	numTweets, err := api.parseNumTweets(c)
	if err != nil {
		api.sendError(c, err)
		return
	}

	if makeFake, ok := c.GetQuery("fake"); ok {
		if makeFake == "true" {
			api.fake(c, ord, numTweets)
			return
		}
		api.real(c, numTweets)
		return
	}

	api.both(c, ord, numTweets)
}

func (api *TrumpAPI) parseOrd(c *gin.Context) (int, error) {
	ord := 2
	if ordQry, ok := c.GetQuery("ord"); ok {
		if parsed, err := strconv.Atoi(ordQry); err == nil {
			ord = parsed
		}
	}

	if ord > api.config.Markov.MaxOrder {
		return 0, ErrOrderTooHigh
	}
	return ord, nil
}

func (api *TrumpAPI) parseNumTweets(c *gin.Context) (int, error) {
	numTweets := 1
	if num, ok := c.GetQuery("n"); ok {
		if parsed, err := strconv.Atoi(num); err == nil {
			numTweets = parsed
		}
	}

	if numTweets < 1 || numTweets > api.config.MaxTweets {
		return 0, ErrNumTweets
	}
	return numTweets, nil
}

func (api *TrumpAPI) fake(c *gin.Context, markovOrder, numTweets int) {
	api.createTweets(c, numTweets, func() ([]tt.Tweet, error) {
		tweet, err := tt.RandomFakeSample(markovOrder, api.config)
		if err != nil {
			return nil, err
		}
		return []tt.Tweet{tweet}, nil
	})
}

func (api *TrumpAPI) real(c *gin.Context, num int) {
	gen := func() ([]tt.Tweet, error) {
		tweet, err := tt.RandomRealSample(&api.config.Markov)
		if err != nil {
			return nil, err
		}
		return []tt.Tweet{tweet}, nil
	}

	api.createTweets(c, num, gen)
}

func (api *TrumpAPI) both(c *gin.Context, order, numTweets int) {
	// Implementation for both real and fake tweets
}

func (api *TrumpAPI) createTweets(c *gin.Context, num int, tweetGen func() ([]tt.Tweet, error)) {
	var tweets []tt.Tweet
	for i := 0; i < num; i++ {
		genTweets, err := tweetGen()
		if err != nil {
			api.sendError(c, err)
			return
		}
		tweets = append(tweets, genTweets...)
	}
	c.JSON(200, tweets)
}

func (api *TrumpAPI) sendError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}
