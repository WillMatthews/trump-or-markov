package trumptweets

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/bcicen/jstream"
)

func LoadTrumpTweets(cfg config.Dataset) {
	jsonFile := cfg.Trump
	f, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := jstream.NewDecoder(f, 1)
	err = parseStream(decoder)
	if err != nil {
		panic(err)
	}
}

func parseStream(decoder *jstream.Decoder) error {
	numTweets := 0
	for mv := range decoder.Stream() {
		value := mv.Value
		switch value.(type) {
		case map[string]interface{}:
			numTweets++
			tweet, err := parseTweet(value.(map[string]interface{}))
			if err != nil {
				// TODO log as error
				return fmt.Errorf("Error parsing tweet: %w", err)
			}
			storeTweet(tweet)

		default:
			return fmt.Errorf("Unexpected type: %T", value)
		}
	}

	fmt.Printf("Loaded %d tweets\n", numTweets)
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

// RandomRealSample returns a random tweet from the dataset
func RandomRealSample(cfg *config.Markov) (Tweet, error) {
	filters := []TweetFilter{
		MinWordsFilter(cfg.MinWords),
	}

	f := func() (Tweet, error) {
		// TODO: replace with SQLite in future
		sample := rand.IntN(len(tweets))
		return tweets[sample], nil
	}

	return randomSampleWithFilter(
		ComposeFilters(filters...),
		f,
		cfg.MaxGenerateAttempts,
	)
}
