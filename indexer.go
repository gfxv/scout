package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const pathsBufferSize = 200
const numIndexWorkers = 5

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

func (s *SearchQueryResult) Path() string {
	return s.path
}

func (s *SearchQueryResult) Rank() float32 {
	return s.rank
}

type Indexer struct {
	diMu          *sync.Mutex
	dfMu          *sync.Mutex
	documentIndex DocIndex
	docFrequency  TermFreq
}

func NewIndexer() *Indexer {
	return &Indexer{
		diMu:          &sync.Mutex{},
		dfMu:          &sync.Mutex{},
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

func (i *Indexer) IndexDir(path string) error {
	filesChan := make(chan string, pathsBufferSize)
	go func() {
		if err := collectFiles(path, filesChan); err != nil {
			fmt.Printf("error occurred during collecting files, err: %w", err)
		}
	}()

	wg := sync.WaitGroup{}
	for w := 0; w < numIndexWorkers; w++ {
		wg.Add(1)
		go func(i *Indexer) {
			i.indexWorker(&wg, filesChan)
		}(i)
	}
	wg.Wait()

	return nil
}

func (i *Indexer) IndexFile(path string) error {
	fmt.Printf("Indexing %s...\n", path)

	// TODO: support other file types (pdf, html, doc?, ...)
	var content []rune
	var err error
	switch filepath.Ext(path) {
	case ".md", ".txt":
		content, err = plainTextReader(path)
	case ".pdf":
		content, err = pdfReader(path)
	default:
		fmt.Printf("Unknown file type %s\n")
		return nil
	}

	if err != nil {
		return err
	}

	i.diMu.Lock()
	if _, ok := i.documentIndex[path]; ok {
		i.RemoveFile(path)
	}
	i.diMu.Unlock()

	tokenizer := NewTokenizer(content)

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
			i.dfMu.Lock()
			i.docFrequency[token]++
			i.dfMu.Unlock()
		}
	}

	i.diMu.Lock()
	i.documentIndex[path] = docInfo
	i.diMu.Unlock()
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

func (i *Indexer) indexWorker(wg *sync.WaitGroup, paths <-chan string) {
	defer wg.Done()
	for path := range paths {
		if err := i.IndexFile(path); err != nil {
			// TODO: write errors to errChan
			fmt.Println(err)
		}
	}
}

func collectFiles(root string, out chan<- string) error {
	defer close(out)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to collect file %s, err: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}

		out <- path
		return nil
	}
	if err := filepath.Walk(root, walkFunc); err != nil {
		return err
	}
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
