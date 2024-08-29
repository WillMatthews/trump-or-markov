package trumptweets

import (
	"fmt"
	"time"
)

// Hold in memory for now, but we will want a SQLite DB
var tweets []Tweet

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

	IsReal bool `json:"isReal"`
}

func randomSampleWithFilter(
	filter TweetFilter,
	generator func() (Tweet, error),
	maxAttempts int,
) (Tweet, error) {
	for i := 0; i < maxAttempts; i++ {
		tweet, err := generator()
		if err != nil {
			return Tweet{}, err
		}
		if filter(tweet) {
			return tweet, nil
		}
	}
	return Tweet{}, fmt.Errorf("could not find suitable tweet after %d attempts", maxAttempts)
}
