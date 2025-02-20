package engine

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/gfxv/scout/internal/database"
	"github.com/gfxv/scout/internal/models"
)

const pathsBufferSize = 200
const numIndexWorkers = 10

type Indexer struct {
	buffer    []database.DocumentData
	batchSize int
	db        *database.Database

	diMu          *sync.Mutex
	dfMu          *sync.Mutex
	documentIndex models.DocIndex
	docFrequency  models.TermFreq
}

func NewIndexer() *Indexer {

	db, err := database.NewDatabase("meta.db")
	if err != nil {
		panic(err)
	}

	return &Indexer{
		buffer:    make([]database.DocumentData, 0),
		batchSize: 50,
		db:        db,

		diMu:          &sync.Mutex{},
		dfMu:          &sync.Mutex{},
		documentIndex: make(models.DocIndex),
		docFrequency:  make(models.TermFreq),
	}
}

func (i *Indexer) SearchQuery(query string) []models.SearchQueryResult {

	return nil
}

func (i *Indexer) Load() error {
	return nil
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
	i.Flush()
	return nil
}

func (i *Indexer) IndexFile(path string) error {
	fmt.Printf("Indexing %s...\n", path)

	// TODO:
	// - may be change switch-case to map ? (command pattern)
	var content []rune
	var err error
	switch filepath.Ext(path) {
	case ".md", ".txt":
		content, err = plainTextReader(path)
	case ".xml", ".xhtml":
		content, err = xmlReader(path)
	case ".pdf":
		content, err = pdfReader(path)
	case ".html":
		content, err = htmlReader(path)
	default:
		fmt.Printf("Unknown file type %s\n", path)
		return nil
	}

	if err != nil {
		return err
	}

	tokenizer := NewTokenizer(content)
	terms := make([]string, 0)
	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}
		terms = append(terms, token)
	}

	i.buffer = append(i.buffer, database.DocumentData{Path: path, Terms: terms})
	if len(i.buffer) >= i.batchSize {
		if err := i.db.AddDocuments(i.buffer); err != nil {
			fmt.Printf("failed to add documents, err: %v", err)
		}
		i.buffer = nil
	}
	return nil
}

// Flush processes any remaining documents in the buffer
func (i *Indexer) Flush() {
	if len(i.buffer) > 0 {
		if err := i.db.AddDocuments(i.buffer); err != nil {
			fmt.Printf("failed to add documents, err: %v", err)
		}
		i.buffer = nil
	}
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

func (i *Indexer) tf(token string, docInfo models.DocInfo) float32 {
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

func tokenizeQuery(query string) []string {
	tokenizer := NewTokenizer([]rune(query))
	tokens := make([]string, 0)
	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func (i *Indexer) prettyPrint() {
	for path, docInfo := range i.documentIndex {
		fmt.Printf("%s (%d total terms)\n", path, docInfo.TotalTerms)
		for term, freq := range docInfo.Terms {
			fmt.Printf("  %s -> %d\n", term, freq)
		}
	}
}
