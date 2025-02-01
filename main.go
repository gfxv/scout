package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

const dirPath = "files"
const sampleFile = "files/scylla-readme.md"

func main() {
	indexer := NewIndexer()
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		path := filepath.Join(dirPath, file.Name())
		fmt.Printf("Indexing %s...\n", path)
		if err := indexer.IndexFile(path); err != nil {
			fmt.Println(err)
		}
	}

	app := fiber.New()

	app.Get("/", adaptor.HTTPHandler(templ.Handler(IndexPage())))
	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		fmt.Println("Received query:", query)
		results := indexer.SearchQuery(query)

		handler := adaptor.HTTPHandler(templ.Handler(searchResults(results)))
		return handler(c)
	})

	if err = app.Listen(":6969"); err != nil {
		fmt.Println(err)
	}

}
