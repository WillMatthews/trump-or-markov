package markov

import (
	"math/rand/v2"
	"strings"
)

type Chain struct {
	chain map[string][]string
	seeds []string
	order int

	stopOnStopProbabilty float64
}

func NewMarkovChain(order int) *Chain {
	return &Chain{
		chain:                make(map[string][]string),
		order:                order,
		stopOnStopProbabilty: 0.6,
	}
}

func (c *Chain) makeKey(words wordChain) string {
	if len(words) == 0 {
		panic("makeKey called with no words")
	}

	// making words lowercase can help with the chain
	// there are more possible paths it can take
	lWords := make([]string, len(words))
	for i, word := range words {
		lWords[i] = strings.ToLower(word)
	}

	key := lWords[0]
	if len(lWords) == 1 {
		return key
	}

	end := min(c.order, len(lWords))
	for _, word := range lWords[1:end] {
		key += " " + word
	}
	return key
}

func (c *Chain) Train(words []string) {
	addEntry := func(key string, value string) {
		if _, ok := c.chain[key]; !ok {
			c.chain[key] = make([]string, 0)
		}
		c.chain[key] = append(c.chain[key], value)
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
	return c.Generate(seed, length)
}

func (c *Chain) Generate(seed string, length int) string {
	words := newWordChain()
	words.add(seed)

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

		next := posible[rand.IntN(len(posible))]
		words.add(next)

		if c.shouldIStop(words) {
			break
		}
	}
	return words.String()
}

func (c *Chain) shouldIStop(words wordChain) bool {
	next := words[len(words)-1]

	if next[len(next)-1] == '.' {
		if rand.Float64() < c.stopOnStopProbabilty {
			return true
		}
	}

	return words.len() > words.len()
}

type wordChain []string

func (w wordChain) String() string {
	doubleSpaceProb := 0.05

	sb := strings.Builder{}
	for i, word := range w {
		if i == 0 {
			sb.WriteString(word)
		} else {
			if rand.Float64() < doubleSpaceProb {
				sb.WriteString(" ")
			}
			sb.WriteString(" ")
			sb.WriteString(word)
		}
	}
	return sb.String()
}

func newWordChain() wordChain {
	return make([]string, 0)
}

func (w wordChain) len() int {
	total := 0
	for _, word := range w {
		total += len(word)
	}
	return total
}

func (w *wordChain) add(word string) {
	*w = append(*w, word)
}
