package main

import (
	"bytes"
	"fmt"
	"os"
)

type TermFreq map[string]uint
type DocIndex map[string]TermFreq

type Indexer struct {
	documentIndex DocIndex
	globalInfo    TermFreq
}

func NewIndexer() *Indexer {
	return &Indexer{
		documentIndex: make(DocIndex),
		globalInfo:    make(TermFreq),
	}
}

func (i *Indexer) IndexFile(path string) error {
	rawFile, err := os.ReadFile(sampleFile)
	if err != nil {
		return fmt.Errorf("can't read file %s, err: %w\n", rawFile, err)
	}
	tokenizer := NewTokenizer(bytes.Runes(rawFile))

	i.documentIndex[path] = make(TermFreq)

	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}

		// add token to global frequencies table
		globalFreq, ok := i.globalInfo[token]
		if ok {
			i.globalInfo[token] = globalFreq + 1
		} else {
			i.globalInfo[token] = 1
		}

		// add token to document frequencies table
		terms := i.documentIndex[path]
		freq, ok := terms[token]
		if ok {
			terms[token] = freq + 1
		} else {
			terms[token] = 1
		}
	}
	return nil
}

func (i *Indexer) RemoveFile(path string) error {
	terms, ok := i.documentIndex[path]
	if !ok {
		return fmt.Errorf("path %s not found", path)
	}
	for term, freq := range terms {
		globalTermFreq, ok := i.globalInfo[term]
		if !ok {
			continue
		}
		i.globalInfo[term] = globalTermFreq - freq
	}
	delete(i.documentIndex, path)
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
