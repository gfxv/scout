package main

import (
	"fmt"
)

const dirPath = "files"
const sampleFile = "files/scylla-readme.md"

func main() {
	fmt.Println("Hello world!")

	indexer := NewIndexer()
	if err := indexer.IndexFile(sampleFile); err != nil {
		fmt.Println(err)
	}

	indexer.prettyPrint()
}
