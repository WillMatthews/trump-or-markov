package trumptweets

import "strings"

type TweetFilter func(*Tweet) bool

func ComposeFilters(filters ...TweetFilter) TweetFilter {
	return func(t *Tweet) bool {
		for _, filter := range filters {
			if !filter(t) {
				return false
			}
		}
		return true
	}
}

func MinWordsFilter(minWords int) TweetFilter {
	return func(t *Tweet) bool {
		return len(strings.Fields(t.Text)) >= minWords
	}
}

// NoEllipsisFilter returns a filter that returns true if the tweet does not contain
// In the training data, ellipsis are used to indicate that the tweet was cut off,
// this can happen if the tweet is too long, is a twitlonger link, or if the tweet
// is a retweet.
func NoEllipsisFilter() TweetFilter {
	return func(t *Tweet) bool {
		return !strings.Contains(t.Text, "…") && !strings.Contains(t.Text, "...")
	}
}
