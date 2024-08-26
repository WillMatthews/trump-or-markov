package trumptweets

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/markov"
)

const (
	maxOrder = 4
	minOrder = 1

	doubleSpaceProb = 0.05
)

var chains map[int]*markov.Chain

func TrainMarkovChain(order int) *markov.Chain {
	chain := markov.NewMarkovChain(order)

	trained := 0
	for _, tweet := range tweets {
		tokens := markov.Tokenise(tweet.Text)
		chain.Train(tokens)
		trained++
	}
	fmt.Printf("Trained on %d tweets\n", trained)
	return chain
}

func getChain(order int) (markov.Chain, error) {
	if chains == nil {
		chains = make(map[int]*markov.Chain)
	}
	chain, ok := chains[order]
	if !ok {
		chain = TrainMarkovChain(order)
		chains[order] = chain
	}

	// j := 0
	// for i, c := range *&chain.Chain {
	// 	fmt.Println(i, c)
	// 	j++
	// 	if j > 20 {
	// 		break
	// 	}
	// }

	return *chain, nil
}

func RandomFakeSample(order int) (*Tweet, error) {
	if order < minOrder || order > maxOrder {
		return nil, fmt.Errorf("order must be between %d and %d", minOrder, maxOrder)
	}

	baseTweet, err := withNoEllipsis(RandomSample)
	if err != nil {
		return nil, err
	}

	chain, err := getChain(order)
	if err != nil {
		return nil, err
	}

	generated := chain.GenerateRandom(order, 140)
	baseTweet.Text = randomSpaceInjection(generated)

	baseTweet.IsRetweet = strings.Contains(baseTweet.Text, "RT")
	return baseTweet, nil
}

func randomSpaceInjection(text string) string {
	words := strings.Fields(text)
	for i, _ := range words {
		if rand.Float64() < doubleSpaceProb {
			words[i] += " "
		}
	}
	return strings.Join(words, " ")
}

// withNoEllipsis returns a tweet that does not contain an ellipsis
// In the training data, ellipsis are used to indicate that the tweet was cut off,
// this can happen if the tweet is too long, is a twitlonger link, or if the tweet
// is a retweet.
func withNoEllipsis(f func() (*Tweet, error)) (*Tweet, error) {
	var sampled *Tweet
	var err error
	for sampled == nil || strings.Contains(sampled.Text, "…") {
		sampled, err = f()
		if err != nil {
			return nil, err
		}
	}
	return sampled, err
}
