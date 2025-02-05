package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

func plainTextReader(path string) ([]rune, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s, err: %w", path, err)
	}
	return bytes.Runes(content), nil
}

func pdfReader(path string) ([]rune, error) {
	file, reader, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s, err: %w", path, err)
	}
	defer file.Close()

	// TODO: figure out how to estimate sclice len(cap) based on file size ???
	content := make([]rune, 0)
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		p := reader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		result, err := p.GetPlainText(nil)
		if err != nil {
			fmt.Println("failed to get plain text from page %d from document %s, err: %w", pageIndex, path, err)
			continue
		}
		content = append(content, []rune(result)...)
	}
	return content, nil
}

// TODO: implement...
func htmlReader(path string) ([]rune, error) {
	return nil, nil
}
