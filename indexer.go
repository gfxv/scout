package main

import (
	"bytes"
	"fmt"
	"os"
)

type TermInfo map[string]uint
type DocIndex map[string]TermInfo

type Indexer struct {
	documentIndex DocIndex
}

func NewIndexer() *Indexer {
	return &Indexer{documentIndex: make(DocIndex)}
}

func (i *Indexer) IndexFile(path string) error {
	rawFile, err := os.ReadFile(sampleFile)
	if err != nil {
		return fmt.Errorf("can't read file %s, err: %w\n", rawFile, err)
	}
	tokenizer := NewTokenizer(bytes.Runes(rawFile))

	// check if file (path) was already indexed
	// create new entry if not
	if _, ok := i.documentIndex[path]; !ok {
		i.documentIndex[path] = make(TermInfo)
	}

	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}
		terms := i.documentIndex[path]

		freq, ok := terms[token]
		if !ok {
			terms[token] = 1
			continue
		}
		terms[token] = freq + 1
	}
	return nil
}

// helper function to print term frequencies
func (i *Indexer) prettyPrint() {
	for path, termInfo := range i.documentIndex {
		fmt.Println(path)
		for term, freq := range termInfo {
			fmt.Printf("  %s -> %d\n", term, freq)
		}
	}
}
