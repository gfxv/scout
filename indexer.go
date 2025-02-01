package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
)

type TermFreq map[string]uint
type DocInfo struct {
	Terms      TermFreq
	TotalTerms uint
}
type DocIndex map[string]DocInfo

type SearchQueryResult struct {
	path string
	rank float32
}

type Indexer struct {
	documentIndex DocIndex
	docFrequency  TermFreq
}

func NewIndexer() *Indexer {
	return &Indexer{
		documentIndex: make(DocIndex),
		docFrequency:  make(TermFreq),
	}
}

func (i *Indexer) SearchQuery(query string) []SearchQueryResult {
	result := make([]SearchQueryResult, 0)
	tokenizer := NewTokenizer([]rune(query))
	tokens := make([]string, 0)

	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}

	for path, docInfo := range i.documentIndex {
		rank := float32(0)
		for _, token := range tokens {
			tf := i.tf(token, docInfo)
			idf := i.idf(token)
			rank += tf * idf
		}
		result = append(result, SearchQueryResult{path, rank})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].rank > result[j].rank
	})

	return result
}

func (i *Indexer) IndexFile(path string) error {
	if _, ok := i.documentIndex[path]; ok {
		i.RemoveFile(path)
	}

	rawFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read file %s, err: %w", path, err)
	}
	tokenizer := NewTokenizer(bytes.Runes(rawFile))

	docInfo := DocInfo{
		Terms:      make(TermFreq),
		TotalTerms: 0,
	}

	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}

		docInfo.Terms[token]++
		docInfo.TotalTerms++

		if docInfo.Terms[token] == 1 {
			i.docFrequency[token]++
		}
	}

	i.documentIndex[path] = docInfo
	return nil
}

func (i *Indexer) RemoveFile(path string) error {
	docInfo, ok := i.documentIndex[path]
	if !ok {
		return fmt.Errorf("path %s not found", path)
	}

	for term := range docInfo.Terms {
		if i.docFrequency[term] > 0 {
			i.docFrequency[term]--
		}
		if i.docFrequency[term] == 0 {
			delete(i.docFrequency, term)
		}
	}

	delete(i.documentIndex, path)
	return nil
}

func (i *Indexer) tf(token string, docInfo DocInfo) float32 {
	if docInfo.TotalTerms == 0 {
		return 0
	}
	freq := docInfo.Terms[token]
	return float32(freq) / float32(docInfo.TotalTerms)
}

func (i *Indexer) idf(token string) float32 {
	totalDocs := len(i.documentIndex)
	if totalDocs == 0 {
		totalDocs = 1
	}
	docsWithTerm := i.docFrequency[token]
	return float32(math.Log(float64(totalDocs+1) / float64(docsWithTerm+1)))
}

func (i *Indexer) prettyPrint() {
	for path, docInfo := range i.documentIndex {
		fmt.Printf("%s (%d total terms)\n", path, docInfo.TotalTerms)
		for term, freq := range docInfo.Terms {
			fmt.Printf("  %s -> %d\n", term, freq)
		}
	}
}
