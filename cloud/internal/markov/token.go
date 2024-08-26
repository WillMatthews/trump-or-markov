package markov

import (
	"strings"
)

type token string

func NewToken(s string) *token {
	interned := dict.Intern(token(s))
	return interned
}

func (t token) String() string {
	return string(t)
}

func (t *token) Lower() token {
	return token(strings.ToLower(t.String()))
}

// Keep this simple for now.
func Tokenise(s string) []*token {
	var tokens []*token

	words := strings.Fields(s)
	for _, word := range words {
		tok := NewToken(word) // Interned
		tokens = append(tokens, tok)
	}
	return tokens
}

func (t token) Len() int {
	return len(t)
}

type tokenChain []*token

func (w tokenChain) String() string {
	var sb strings.Builder
	for i, word := range w {
		if i == 0 {
			sb.WriteString(word.String())
		} else {
			sb.WriteString(" ")
			sb.WriteString(word.String())
		}
	}
	return sb.String()
}

func NewTokenChain() tokenChain {
	return make([]*token, 0)
}

func (w tokenChain) Len() int {
	total := 0
	for _, word := range w {
		total += word.Len()
	}
	return total
}

func (w *tokenChain) Add(tok *token) {
	*w = append(*w, tok)
}
