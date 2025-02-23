package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// TODO: add .env parsing

type Config struct {
	Index  bool   `env:"SCOUT_INDEX" envDefault:"false"`
	Serve  bool   `env:"SCOUT_SERVE" envDefault:"false"`
	Files  string `env:"SCOUT_FILES"`
	Port   string `env:"SCOUT_PORT" envDefault:"6969"`
	DBPath string `env:"SCOUT_DB_PATH" envDefault:"meta.db"`
}

func ParseConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Failed to load .env file: %v", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Warnf("Failed to parse environment variables: %v\n", err)
	}

	fmt.Printf("%v+\n", cfg)

	flag.BoolVar(&cfg.Index, "index", cfg.Index, "Index files in the specified directory and exit (required with -files)")
	flag.BoolVar(&cfg.Serve, "serve", cfg.Serve, "Start the search web server")
	flag.StringVar(&cfg.Files, "files", cfg.Files, "Directory path containing files to index (required with -index)")
	flag.StringVar(&cfg.Port, "port", cfg.Port, "Port to listen on when serving (e.g., 8080)")
	flag.StringVar(&cfg.DBPath, "db", cfg.DBPath, "Path to the SQLite database file")
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
