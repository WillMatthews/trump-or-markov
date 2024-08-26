package markov

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/WillMatthews/trump-or-markov/internal/config"
)

type Chain struct {
	Chain map[state]stateTransitions
	seeds []token
	order int

	endPunctuation       []string
	stopOnStopProbabilty float64
}

type stateTransitions []nextState

func (states stateTransitions) sample() *nextState {
	if len(states) == 0 {
		return nil
	}

	var mass uint16
	for _, f := range states {
		mass += f.ProbMass
	}

	if mass == 0 {
		// This should never happen
		// TODO log an error here
		return states.unifSample()
	}

	r := rand.IntN(int(mass))
	for _, f := range states {
		r -= int(f.ProbMass)
		if r <= 0 {
			return &f
		}
	}

	// This should never happen
	// TODO log an error here
	return states.unifSample()
}

func (states stateTransitions) unifSample() *nextState {
	if len(states) == 0 {
		return nil
	}

	return &states[rand.IntN(len(states))]
}

type nextState struct {
	TokenIdx tokenIndex
	ProbMass uint16
}

type state string

func (k state) String() string {
	return string(k)
}

func NewMarkovChain(
	order int,
	config *config.Markov,
) *Chain {
	return &Chain{
		Chain:                make(map[state]stateTransitions),
		order:                order,
		stopOnStopProbabilty: config.EndPunctuationProb,
		endPunctuation:       config.EndPunctuation,
	}
}

func (c *Chain) makeKey(tokens tokenChain) state {
	if len(tokens) == 0 {
		panic("makeKey called with no words")
	}

	// making words lowercase can help with the chain in our case.
	// this reduces the number of states in the chain and makes
	// it a bit less sparse, which can help avoid overfitting.
	var nTokens = make(tokenChain, len(tokens))
	for i, t := range tokens {
		nTokens[i] = getOrCreateTokenIndex(t.Token())
	}

	end := len(nTokens) - c.order
	if end <= 0 {
		return GetState(nTokens)
	}
	return GetState(nTokens[:end])
}

func incrementCount(freq stateTransitions, value token) error {
	for i, f := range freq {
		if f.TokenIdx.Token() == value {
			freq[i].ProbMass++
			return nil
		}
	}
	return errors.New("value not found")
}

func (c *Chain) addEntry(key state, value token) {
	tIdx := getOrCreateTokenIndex(value)

	// hack hack hack
	// ensure lowercase token is stored (which is used in makeKey)
	_ = getOrCreateTokenIndex(value.Lower())

	if freq, ok := c.Chain[key]; !ok {
		c.Chain[key] = append(
			c.Chain[key],
			nextState{TokenIdx: tIdx, ProbMass: 1},
		)
	} else {
		if err := incrementCount(freq, value); err != nil {
			newVal := nextState{TokenIdx: tIdx, ProbMass: 1}
			c.Chain[key] = append(c.Chain[key], newVal)
		}
	}
}

func (c *Chain) Train(words []token) {
	if len(words) == 0 {
		return
	}

	c.seeds = append(c.seeds, words[0])
	for i := 1; i < len(words); i++ {
		start := 0
		if i > c.order {
			start = i - c.order
		}
		tChain := tokenSliceToChain(words[start:i])
		key := c.makeKey(tChain)
		c.addEntry(key, words[i])
	}
}

func tokenSliceToChain(words []token) tokenChain {
	chain := make(tokenChain, len(words))
	for i, w := range words {
		chain[i] = getOrCreateTokenIndex(w)
	}
	return chain
}

func (c *Chain) GenerateRandom(order, length int) string {
	// MUST retry with new seed if failure
	// otherwise a high probability of getting stuck
	seedIdx := rand.IntN(len(c.seeds))
	seed := c.seeds[seedIdx]
	genTokens := c.Generate(seed, length)

	var sb strings.Builder
	for _, t := range genTokens {
		sb.WriteString(t.Token().String())
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

		possible, ok := c.Chain[key]
		if !ok {
			break
		}

		// sample the possible next state
		next := possible.sample()
		if next == nil {
			break
		}
		words.Add(next.TokenIdx.Token())

		if c.decideStop(words, length) {
			break
		}
	}
	return words
}

func (c *Chain) decideStop(words tokenChain,
	stopWordLimit int,
) bool {
	next := words[len(words)-1]
	nextStr := next.Token().String()
	endChar := string(nextStr[len(nextStr)-1])

	// hmm - do as tokens instead?
	if slices.Contains(c.endPunctuation, endChar) {
		if rand.Float64() < c.stopOnStopProbabilty {
			return true
		}
	}

	return words.Len() > stopWordLimit
}

func GetState(words tokenChain) state {
	if len(words) == 0 {
		panic("NewKey called with no words")
	}

	pruned := pruneWordsToOrder(words, 2)

	keywords := pruned[0].Token().String()
	if len(pruned) == 1 {
		return state(keywords)
	}

	for _, word := range pruned[1:] {
		keywords += " " + word.Token().String()
	}

	return state(keywords)
}

func pruneWordsToOrder(words tokenChain, order int) tokenChain {
	if len(words) <= order {
		return words
	}

	return words[len(words)-order:]
}
