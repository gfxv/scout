package main

import (
	"bytes"
	"fmt"
	"os"
)

const dirPath = "files"
const sampleFile = "files/scylla-readme.md"

func main() {
	fmt.Println("Hello world!")

	rawFile, err := os.ReadFile(sampleFile)
	if err != nil {
		fmt.Printf("can't read file %s, err: %s\n", rawFile, err.Error())
		return
	}

	tokenizer := NewTokenizer(bytes.Runes(rawFile))
	for {
		token, ok := tokenizer.NextToken()
		if !ok {
			break
		}
		fmt.Println(token)
	}
}
