package markov

import (
	"strings"
)

type tokenIndex int

var (
	globalTokenPool    = make(map[token]tokenIndex)
	globalTokenPoolRev = make(map[tokenIndex]token)
	nextTokenIndex     tokenIndex
)

func getTokenIndex(s token) tokenIndex {
	if idx, exists := globalTokenPool[s]; exists {
		return idx
	}
	return -1
}

func getOrCreateTokenIndex(s token) tokenIndex {
	if idx, exists := globalTokenPool[s]; exists {
		return idx
	}

	idx := nextTokenIndex
	globalTokenPool[s] = idx
	globalTokenPoolRev[idx] = s
	nextTokenIndex++
	return idx
}

func (t tokenIndex) Token() token {
	return globalTokenPoolRev[t]
}

type token string

func (t token) Index() tokenIndex {
	return getOrCreateTokenIndex(t)
}

func (t token) String() string {
	return string(t)
}

func (t *token) Lower() token {
	return token(strings.ToLower(t.String()))
}

func Tokenise(s string) []token {
	words := strings.Fields(s)
	var tokens []token
	for _, word := range words {
		tok := token(word)
		tokens = append(tokens, tok)
	}
	return tokens
}

func (t token) Len() int {
	return len(t)
}

type tokenChain []tokenIndex

func (w tokenChain) String() string {
	var sb strings.Builder
	for i, word := range w {
		if i == 0 {
			sb.WriteString(word.Token().String())
		} else {
			sb.WriteString(" ")
			sb.WriteString(word.Token().String())
		}
	}
	return sb.String()
}

func NewTokenChain() tokenChain {
	return make([]tokenIndex, 0)
}

func (w tokenChain) Len() int {
	total := 0
	for _, word := range w {
		total += word.Token().Len()
	}
	return total
}

func (w *tokenChain) Add(tok token) {
	*w = append(*w, tok.Index())
}
