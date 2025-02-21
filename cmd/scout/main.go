package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gfxv/scout/internal/config"
	"github.com/gfxv/scout/internal/engine"
	"github.com/gfxv/scout/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	cfg := config.ParseConfig()
	if err := config.ValidateConfig(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	indexer := engine.NewIndexer()
	var err error
	switch {
	case cfg.Index:
		err = runIndexer(indexer, cfg.Files)
	case cfg.Serve:
		err = runServer(indexer, cfg.Port)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

}

func runIndexer(indexer *engine.Indexer, dir string) error {
	start := time.Now()
	if err := indexer.IndexDir(dir); err != nil {
		return fmt.Errorf("indexing failed: %v", err)
	}
	fmt.Printf("Indexing took: %s\n", time.Since(start))
	return nil
}

func runServer(indexer *engine.Indexer, port string) error {
	if err := indexer.Load(); err != nil {
		return fmt.Errorf("failed to load indexer: %v", err)
	}

	app := fiber.New()
	app.Get("/", adaptor.HTTPHandler(templ.Handler(views.IndexPage())))
	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		log.Printf("Received query: %s\n", query)
		results := indexer.SearchQuery(query)
		handler := adaptor.HTTPHandler(templ.Handler(views.SearchResults(results)))
		return handler(c)
	})

	listenAddr := ":" + port
	if err := app.Listen(listenAddr); err != nil {
		return fmt.Errorf("server failed on %s: %v", listenAddr, err)
	}
	return nil
}
