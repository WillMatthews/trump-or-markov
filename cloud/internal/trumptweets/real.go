package trumptweets

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/bcicen/jstream"
	"github.com/rs/zerolog/log"
)

// Hold in memory for now, but we will want a SQLite DB
var tweets []Tweet

type Tweet struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Favorites int       `json:"favorites"`
	Retweets  int       `json:"retweets"`
	Date      time.Time `json:"date"`
	Device    string    `json:"device"`

	IsRetweet bool `json:"isRetweet"`
	IsDeleted bool `json:"isDeleted"`

	IsFlagged bool `json:"isFlagged"`
}

type DirtyTweet struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Favorites int    `json:"favorites"`
	Retweets  int    `json:"retweets"`
	Date      string `json:"date"`
	Device    string `json:"device"`

	IsRetweet string `json:"isRetweet"`
	IsDeleted string `json:"isDeleted"`

	IsFlagged string `json:"isFlagged"`
}

func LoadTrumpTweets(cfg config.Dataset) {
	jsonFile := cfg.Trump
	f, err := os.Open(jsonFile)
	if err != nil {
		log.Fatal().Err(err).
			Msg("Error opening file")
	}
	defer f.Close()

	decoder := jstream.NewDecoder(f, 1)
	err = parseStream(decoder)
	if err != nil {
		log.Fatal().Err(err).
			Msg("Error parsing stream")
	}
}

func parseStream(decoder *jstream.Decoder) error {
	numTweets := 0
	for mv := range decoder.Stream() {
		value := mv.Value
		switch value.(type) {
		case map[string]interface{}:
			numTweets++
			_, err := parseTweet(value.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("Error parsing tweet: %w", err)
			}

		default:
			return fmt.Errorf("Unexpected type: %T", value)
		}
	}

	log.Info().Int("numTweets", numTweets).
		Msg("Loaded tweets into memory")
	if len(tweets) == 0 {
		return fmt.Errorf("No tweets loaded")
	}
	return nil
}

// TODO replace with SQLite in future
func storeTweet(tweet *Tweet) error {
	tweets = append(tweets, *tweet)
	return nil
}

// RandomSample returns a random tweet from the dataset
// TODO: replace with SQLite in future
func RandomSample() (*Tweet, error) {
	sample := rand.IntN(len(tweets))
	return &tweets[sample], nil
}
