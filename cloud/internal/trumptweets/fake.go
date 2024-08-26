package trumptweets

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/config"
	"github.com/WillMatthews/trump-or-markov/internal/markov"
)

var chains map[int]*markov.Chain

func MakeTweetsChain(order int, cfg *config.Markov) *markov.Chain {
	chain := markov.NewMarkovChain(order, cfg)

	trained := 0
	for _, tweet := range tweets {
		tokens := markov.Tokenise(tweet.Text)
		chain.Train(tokens)
	}
	fmt.Printf("Trained on %d tweets\n", trained)
	return chain
}

func getChain(order int,
	cfg *config.Markov,
) (*markov.Chain, error) {
	if chains == nil {
		chains = make(map[int]*markov.Chain)
	}
	chain, ok := chains[order]
	if !ok {
		chain = MakeTweetsChain(order, cfg)
		chains[order] = chain
	}

	return chain, nil
}

func RandomFakeSample(
	order int,
	config *config.TrumpTwitter,
) (Tweet, error) {
	filters := []TweetFilter{
		MinWordsFilter(config.Markov.MinWords),
		NoEllipsisFilter(),
	}

	generator := func() (Tweet, error) {
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
) (Tweet, error) {
	baseTweet, err := RandomRealSample(&cfg.Markov)
	if err != nil {
		return Tweet{}, err
	}

	chain, err := getChain(order, &cfg.Markov)
	if err != nil {
		return Tweet{}, err
	}

	generated := chain.GenerateRandom(order, cfg.Markov.MaxChars)
	baseTweet.Text = randomSpaceInjection(generated, cfg.DoubleSpaceProb)

	baseTweet.IsRetweet = strings.Contains(baseTweet.Text, "RT")
	baseTweet.IsReal = false
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
