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
	stopChars = []byte{'!', '?', '.'} //  …  (ellipsis is multi-byte, ignore for now)
)

type frequencies []frequency

type frequency struct {
	Token *Token
	Count int
}

type Chain struct {
	chain map[key]frequencies
	seeds []Token
	order int

	stopOnStopProbabilty float64
}

type key string

func hash(s string) string {
	return s
	// h := fnv.New32a()
	// h.Write([]byte(s))
	// return h.Sum32()
}

func NewKey(words []Token) key {
	if len(words) == 0 {
		panic("NewKey called with no words")
	}

	keywords := words[0].String()
	if len(words) == 1 {
		return key(hash(keywords))
	}

	for _, word := range words[1:] {
		keywords += " " + word.String()
	}

	return key(hash(keywords))
}

func (k key) String() string {
	return string(k)
}

func NewMarkovChain(order int) *Chain {
	return &Chain{
		chain:                make(map[key]frequencies),
		order:                order,
		stopOnStopProbabilty: haltOnStopProbabilty,
	}
}

func (c *Chain) makeKey(tokens tokenChain) key {
	if len(tokens) == 0 {
		panic("makeKey called with no words")
	}

	// making words lowercase can help with the chain
	// there are more possible paths it can take
	nTokens := make([]Token, len(tokens))
	for i, t := range tokens {
		nTokens[i] = t.Lower()
	}

	end := len(nTokens) - c.order
	if end <= 0 {
		return NewKey(nTokens)
	}
	return NewKey(nTokens[:end])
}

func (c *Chain) Train(words []Token) {
	incrementCount := func(freq frequencies, value Token) error {
		for i, f := range freq {
			if *f.Token == value {
				freq[i].Count++
				return nil
			}
		}
		return errors.New("value not found")
	}

	addEntry := func(key key, value Token) {
		if freq, ok := c.chain[key]; !ok {
			c.chain[key] = append(c.chain[key], frequency{Token: &value, Count: 1})
		} else {
			if err := incrementCount(freq, value); err != nil {
				newVal := frequency{Token: &value, Count: 1}
				c.chain[key] = append(c.chain[key], newVal)
			}
		}
	}

	// Seeds!
	c.seeds = append(c.seeds, words[0])
	// TODO relace with range over int
	for i := 1; i < c.order+1; i++ {
		if len(words) <= i {
			break
		}
		key := c.makeKey(words[:i])
		addEntry(key, words[i])
	}

	// TODO relace with range over int
	for i := 0; i < len(words)-c.order; i++ {
		key := c.makeKey(words[i : i+c.order])

		addEntry(key, words[i+c.order])
	}
}

func (c *Chain) GenerateRandom(order int, length int) string {
	seedIdx := rand.IntN(len(c.seeds))
	seed := c.seeds[seedIdx]
	genTokens := c.Generate(seed, length)

	var sb strings.Builder
	for _, t := range genTokens {
		sb.WriteString(t.String())
		sb.WriteString(" ")
	}
	return sb.String()
}

func (c *Chain) Generate(seed Token, length int) tokenChain {
	words := NewTokenChain()
	words.Add(seed)

	for {
		start := len(words) - c.order
		if start < 0 {
			start = 0
		}
		wordsForKey := words[start:]
		key := c.makeKey(wordsForKey)

		posible, ok := c.chain[key]
		if !ok {
			break
		}
		next := samplePossibles(posible)
		words.Add(*next.Token)

		if c.shouldIStop(words, length) {
			break
		}
	}
	return words
}

func samplePossibles(possibles frequencies) *frequency {
	mass := 0
	for _, f := range possibles {
		mass += f.Count
	}

	r := rand.IntN(mass)
	for _, f := range possibles {
		r -= f.Count
		if r <= 0 {
			return &f
		}
	}
	return nil
}

func (c *Chain) shouldIStop(words tokenChain, stopWordLimit int) bool {
	next := words[len(words)-1]
	endChar := next[len(next)-1]
	if slices.Contains(stopChars, endChar) {
		if rand.Float64() < c.stopOnStopProbabilty {
			return true
		}
	}

	return words.Len() > stopWordLimit
}
