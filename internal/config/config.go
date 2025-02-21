package config

import (
	"flag"
	"fmt"
)

// TODO: add .env parsing

type Config struct {
	Index  bool
	Serve  bool
	Files  string
	Port   string
	DBPath string
}

func ParseConfig() Config {
	var cfg Config
	flag.BoolVar(&cfg.Index, "index", false, "Index files in the specified directory and exit (required with -files)")
	flag.BoolVar(&cfg.Serve, "serve", false, "Start the search web server")
	flag.StringVar(&cfg.Files, "files", "", "Directory path containing files to index (required with -index)")
	flag.StringVar(&cfg.Port, "port", "6969", "Port to listen on when serving (e.g., 8080)")
	flag.StringVar(&cfg.DBPath, "db", "meta.db", "Path to the SQLite database file")
	flag.Parse()
	return cfg
}

func ValidateConfig(cfg *Config) error {
	if cfg.Index && cfg.Serve {
		return fmt.Errorf("cannot use both -index and -serve together")
	}
	if !cfg.Index && !cfg.Serve {
		return fmt.Errorf("must specify either -index or -serve")
	}
	if cfg.Index && cfg.Files == "" {
		return fmt.Errorf("-files is required when using -index")
	}
	return nil
}
