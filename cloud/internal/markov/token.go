package markov

import (
	"strings"
)

type Token string

func (t Token) String() string {
	return string(t)
}

func (t Token) Lower() Token {
	return Token(strings.ToLower(t.String()))
}

// Keep this simple for now.
func Tokenise(s string) []Token {
	var tokens []Token

	words := strings.Fields(s)
	for _, word := range words {
		tokens = append(tokens, Token(word))
	}
	return tokens
}

func (t Token) Len() int {
	return len(t)
}

type tokenChain []Token

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
	return make([]Token, 0)
}

func (w tokenChain) Len() int {
	total := 0
	for _, word := range w {
		total += len(word)
	}
	return total
}

func (w *tokenChain) Add(tok Token) {
	*w = append(*w, tok)
}
