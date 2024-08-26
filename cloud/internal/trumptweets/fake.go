package trumptweets

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/WillMatthews/trump-or-markov/internal/markov"
)

var chains map[int]*markov.Chain

func TrainMarkovChain(order int) *markov.Chain {
	chain := markov.NewMarkovChain(order, &config.Markov)

	trained := 0
	for _, tweet := range tweets {
		tokens := markov.Tokenise(tweet.Text)
		chain.Train(tokens)
		trained++
	}
	fmt.Printf("Trained on %d tweets\n", trained)
	return chain
}

func getChain(order int) (*markov.Chain, error) {
	if chains == nil {
		chains = make(map[int]*markov.Chain)
	}
	chain, ok := chains[order]
	if !ok {
		chain = TrainMarkovChain(order)
		chains[order] = chain
	}

	return chain, nil
}

func RandomFakeSample(
	order int,
	config *config.TrumpTwitter,
) (*Tweet, error) {
	filters := []TweetFilter{
		MinWordsFilter(config.Markov.MinWords),
		NoEllipsisFilter(),
	}

	generator := func() (*Tweet, error) {
		return generateFake(order, config)
	}

	return randomSampleWithFilter(
		ComposeFilters(filters...),
		generator,
		config.Markov.MaxGenerateAttempts,
	)
}

func generateFake(order int,
	cfg *config.TrumpTwitter,
) (*Tweet, error) {
	baseTweet, err := RandomRealSample(&cfg.Markov)
	if err != nil {
		return nil, err
	}

	chain, err := getChain(order)
	if err != nil {
		return nil, err
	}

	generated := chain.GenerateRandom(order, cfg.Markov.MaxChars)
	baseTweet.Text = randomSpaceInjection(generated, cfg.DoubleSpaceProb)

	baseTweet.IsRetweet = strings.Contains(baseTweet.Text, "RT")
	return baseTweet, nil
}

func randomSpaceInjection(text string,
	doubleSpaceProb float64,
) string {
	words := strings.Fields(text)
	for i, _ := range words {
		if rand.Float64() < doubleSpaceProb {
			words[i] += " "
		}
	}
	return strings.Join(words, " ")
}
