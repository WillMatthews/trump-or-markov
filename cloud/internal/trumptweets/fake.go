package trumptweets

import (
	"fmt"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/markov"
)

var chains map[int]markov.Chain

func TrainMarkovChain(order int) {
	chain := markov.NewMarkovChain(order)

	trained := 0
	for _, tweet := range tweets {
		chain.Train(strings.Fields(tweet.Text))
		trained++
	}
	fmt.Printf("Trained on %d tweets\n", trained)

	chains[order] = *chain
}

func getChain(order int) markov.Chain {
	if chains == nil {
		chains = make(map[int]markov.Chain)
	}
	chain, ok := chains[order]
	if !ok {
		TrainMarkovChain(order)
	}
	chain = chains[order]
	return chain
}

func RandomFakeSample(order int) Tweet {
	base := RandomSample()

	chain := getChain(order)

	base.Text = chain.GenerateRandom(order, 140)

	return base
}
