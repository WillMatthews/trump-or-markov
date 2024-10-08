package trumptweets

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

// We parse the JSON file into a slice of Tweets
func parseDirty(encTweet map[string]interface{}) (*DirtyTweet, error) {
	var decTweet = &DirtyTweet{}
	err := mapstructure.Decode(encTweet, decTweet)
	if err != nil {
		return nil, fmt.Errorf("error decoding tweet: %w", err)
	}
	return decTweet, nil
}

func parseToClean(dirty *DirtyTweet) (*Tweet, error) {
	parsedDate, err := time.Parse("2006-01-02 15:04:05", dirty.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	return &Tweet{
		ID:        dirty.ID,
		Text:      dirty.Text,
		Favorites: dirty.Favorites,
		Retweets:  dirty.Retweets,
		Date:      parsedDate,
		Device:    dirty.Device,

		IsRetweet: dirty.IsRetweet == "t",
		IsDeleted: dirty.IsDeleted == "t",
		IsFlagged: dirty.IsFlagged == "t",
		IsReal:    true,
	}, nil
}

func parseTweet(encTweet map[string]interface{}) (*Tweet, error) {
	dirtyParsed, err := parseDirty(encTweet)
	if err != nil {
		return nil, fmt.Errorf("error parsing tweet: %w", err)
	}

	tweet, err := parseToClean(dirtyParsed)
	if err != nil {
		return nil, fmt.Errorf("error parsing tweet: %w", err)
	}
	return tweet, nil
}
