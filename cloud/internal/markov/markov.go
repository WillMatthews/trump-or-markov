package markov

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strings"
)

const (
	haltOnStopProbabilty = 0.1
)

var (
	stopChars = []byte{'!', '?', '.'} //  …  ellipsis is multi byte. All omitted for now.
)

type Chain struct {
	Chain map[state]stateTransitions
	seeds []token
	order int

	stopOnStopProbabilty float64
}

type stateTransitions []nextState

type nextState struct {
	Token    token // not sure if this does me any good.
	ProbMass int
}

type state string

func (k state) String() string {
	return string(k)
}

func GetState(words []token) state {
	if len(words) == 0 {
		panic("NewKey called with no words")
	}

	pruned := pruneWordsToOrder(words, 2)

	keywords := pruned[0].String()
	if len(pruned) == 1 {
		return state(keywords)
	}

	for _, word := range pruned[1:] {
		keywords += " " + word.String()
	}

	return state(keywords)
}

func pruneWordsToOrder(words []token, order int) []token {
	if len(words) <= order {
		return words
	}

	return words[len(words)-order:]
}

func NewMarkovChain(order int) *Chain {
	return &Chain{
		Chain:                make(map[state]stateTransitions),
		order:                order,
		stopOnStopProbabilty: haltOnStopProbabilty,
	}
}

func (c *Chain) makeKey(tokens tokenChain) state {
	if len(tokens) == 0 {
		panic("makeKey called with no words")
	}

	// making words lowercase can help with the chain
	// there are more possible paths it can take
	nTokens := make([]token, len(tokens))
	for i, t := range tokens {
		nTokens[i] = t.Lower()
	}

	end := len(nTokens) - c.order
	if end <= 0 {
		return GetState(nTokens)
	}
	return GetState(nTokens[:end])
}

func (c *Chain) Train(words []token) {
	incrementCount := func(freq stateTransitions, value token) error {
		for i, f := range freq {
			if f.Token == value {
				freq[i].ProbMass++
				return nil
			}
		}
		return errors.New("value not found")
	}

	addEntry := func(key state, value token) {
		if freq, ok := c.Chain[key]; !ok {
			c.Chain[key] = append(c.Chain[key], nextState{Token: value, ProbMass: 1})
		} else {
			if err := incrementCount(freq, value); err != nil {
				newVal := nextState{Token: value, ProbMass: 1}
				c.Chain[key] = append(c.Chain[key], newVal)
			}
		}
	}

	// Seeds!
	c.seeds = append(c.seeds, words[0])
	// TODO relace with range
	for i := 1; i < c.order+1; i++ {
		if len(words) <= i {
			break
		}
		key := c.makeKey(words[:i])
		addEntry(key, words[i])
	}

	// TODO relace with range
	for i := 0; i < len(words)-c.order; i++ {
		key := c.makeKey(words[i : i+c.order])

		addEntry(key, words[i+c.order])
	}
}

func (c *Chain) GenerateRandom(order, length, minTokens int) string {
	genTokens := withMinTokens(minTokens, func() tokenChain {
		// MUST retry with new seed if failure
		// otherwise a high probability of getting stuck
		seedIdx := rand.IntN(len(c.seeds))
		seed := c.seeds[seedIdx]
		return c.Generate(seed, length)
	})

	var sb strings.Builder
	for _, t := range genTokens {
		sb.WriteString(t.String())
		sb.WriteString(" ")
	}
	return sb.String()
}

func (c *Chain) Generate(seed token, length int) tokenChain {
	words := NewTokenChain()
	words.Add(seed)

	for {
		start := len(words) - c.order
		if start < 0 {
			start = 0
		}
		wordsForKey := words[start:]
		key := c.makeKey(wordsForKey)

		posible, ok := c.Chain[key]
		if !ok {
			break
		}
		next := samplePossibles(posible)
		words.Add(next.Token)

		if c.decideStop(words, length) {
			break
		}
	}
	return words
}

func (c *Chain) decideStop(words tokenChain, stopWordLimit int) bool {
	next := words[len(words)-1]
	endChar := next[len(next)-1]
	// hmm - do as tokens instead?
	if slices.Contains(stopChars, endChar) {
		if rand.Float64() < c.stopOnStopProbabilty {
			return true
		}
	}

	return words.Len() > stopWordLimit
}

func withMinTokens(min int, tweetGen func() tokenChain) tokenChain {
	var words tokenChain
	for len(words) < min {
		words = tweetGen()
	}
	return words
}

func samplePossibles(possibles stateTransitions) *nextState {
	mass := 0
	for _, f := range possibles {
		mass += f.ProbMass
	}

	r := rand.IntN(mass)
	for _, f := range possibles {
		r -= f.ProbMass
		if r <= 0 {
			return &f
		}
	}
	return nil
}
