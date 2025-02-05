package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
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

type TextNode struct {
	Content string `xml:",chardata"`
}

func xmlReader(path string) ([]rune, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s, err: %w", path, err)
	}
	defer file.Close()

	content := make([]rune, 0)
	decoder := xml.NewDecoder(file)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode token from file %s, err: %w", path, err)
		}

		switch t := token.(type) {
		case xml.CharData:
			content = append(content, bytes.Runes(t)...)
		default:
			// ignore other tokens (e.g., start and end tags)
		}
	}
	return content, nil

}

// TODO: implement...
func htmlReader(path string) ([]rune, error) {
	return nil, nil
}
