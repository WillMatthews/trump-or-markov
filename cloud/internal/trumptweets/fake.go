package trumptweets

import (
	"fmt"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/markov"
)

const (
	maxOrder = 4
	minOrder = 1
)

var chains map[int]*markov.Chain

func TrainMarkovChain(order int) *markov.Chain {
	chain := markov.NewMarkovChain(order)

	trained := 0
	for _, tweet := range tweets {
		chain.Train(strings.Fields(tweet.Text))
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

	return *chain, nil
}

func RandomFakeSample(order int) (*Tweet, error) {
	if order < minOrder || order > maxOrder {
		return nil, fmt.Errorf("order must be between %d and %d", minOrder, maxOrder)
	}

	baseTweet, err := RandomSample()
	if err != nil {
		return nil, err
	}

	chain, err := getChain(order)
	if err != nil {
		return nil, err
	}

	baseTweet.Text = chain.GenerateRandom(order, 140)
	baseTweet.IsRetweet = strings.Contains(baseTweet.Text, "RT")
	return baseTweet, nil
}
