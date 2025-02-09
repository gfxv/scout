package main

import (
	"fmt"
	"time"

	"github.com/a-h/templ"
	"github.com/gfxv/scout/internal/engine"
	"github.com/gfxv/scout/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

const dirPath = "files"
const sampleFile = "files/scylla-readme.md"

func main() {
	indexer := engine.NewIndexer()

	start := time.Now()
	if err := indexer.IndexDir(dirPath); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Indexed in: %s", time.Since(start))

	app := fiber.New()

	app.Get("/", adaptor.HTTPHandler(templ.Handler(views.IndexPage())))
	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		fmt.Println("Received query:", query)
		results := indexer.SearchQuery(query)

		handler := adaptor.HTTPHandler(templ.Handler(views.SearchResults(results)))
		return handler(c)
	})

	if err := app.Listen(":6969"); err != nil {
		fmt.Println(err)
	}
}
