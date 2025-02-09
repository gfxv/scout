package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	return &Database{db: nil}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
