package engine

import (
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

type Tokenizer struct {
	data []rune
}

func NewTokenizer(data []rune) *Tokenizer {
	return &Tokenizer{
		data: data,
	}
}

func (t *Tokenizer) NextToken() (string, bool) {
	t.trimLeftSpaces()
	if len(t.data) == 0 {
		return "", false
	}

	if unicode.IsDigit(t.data[0]) {
		token := t.removeIf(func(r rune) bool {
			return unicode.IsDigit(r)
		})
		return string(token), true
	}

	if unicode.IsLetter(t.data[0]) {
		token := t.removeIf(func(r rune) bool {
			return unicode.IsLetter(r) || unicode.IsDigit(r)
		})
		stemmed := t.stem(strings.ToLower(string(token)))
		return stemmed, true
	}

	return string(t.removeFirst(1)), true
}

func (t *Tokenizer) trimLeftSpaces() {
	for len(t.data) > 0 && unicode.IsSpace(t.data[0]) {
		t.data = t.data[1:]
	}
}

func (t *Tokenizer) removeFirst(n int) []rune {
	if n > len(t.data) {
		t.data = []rune{}
		return t.data
	}
	token := t.data[:n]
	t.data = t.data[n:]
	return token
}

func (t *Tokenizer) removeIf(condition func(rune) bool) []rune {
	n := 0
	for n < len(t.data) && condition(t.data[n]) {
		n++
	}
	return t.removeFirst(n)
}

func (t *Tokenizer) stem(token string) string {
	// last boolean argument reduces all the 'useless' words from the language
	// like which, and, why, with, but, etc. (acutally doesn't really matter,
	// because tf-idf handles these words quite good).
	// error may happen only if the language is unknown for the library
	stemmed, _ := snowball.Stem(token, "english", true) // TODO: support different languages
	return stemmed
}
