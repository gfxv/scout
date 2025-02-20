package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gfxv/scout/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DocumentData struct {
	Path  string
	Terms []string
}

type Database struct {
	db *gorm.DB
}

func NewDatabase(path string) (*Database, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory, err: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	/*
		if err := db.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
			return nil, err
		}
		if err := db.Exec("PRAGMA synchronous = NORMAL;").Error; err != nil {
			return nil, err
		}
	*/

	if err := db.AutoMigrate(&models.Document{}, &models.Term{}, &models.DocumentTerm{}); err != nil {
		return nil, fmt.Errorf("failed to create tables, err: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) AddDocuments(docs []DocumentData) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	defer tx.Rollback()

	// Process documents and terms
	documents, termFreqs, uniqueTerms, err := d.processDocuments(docs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Insert documents
	if err := d.insertDocuments(tx, documents); err != nil {
		fmt.Println(err)
		return err
	}

	// Handle terms (insert new, update existing)
	if err := d.handleTerms(tx, uniqueTerms); err != nil {
		fmt.Println(err)
		return err
	}

	// Insert document-term relationships
	if err := d.insertDocumentTerms(tx, documents, termFreqs); err != nil {
		fmt.Println(err)
		return err
	}

	return tx.Commit().Error
}

func (d *Database) LoadIndexData() {}

func (d *Database) processDocuments(docs []DocumentData) ([]models.Document, []map[string]uint, map[string]uint, error) {
	documents := make([]models.Document, 0)
	termFreqs := make([]map[string]uint, len(docs))
	uniqueTerms := make(map[string]uint)

	for i, docData := range docs {
		tf := make(map[string]uint)
		for _, term := range docData.Terms {
			tf[term]++
		}
		totalTerms := uint(len(docData.Terms))
		documents = append(documents, models.Document{Path: docData.Path, TotalTerms: totalTerms})
		termFreqs[i] = tf

		for term := range tf {
			uniqueTerms[term]++
		}
	}

	return documents, termFreqs, uniqueTerms, nil
}

func (d *Database) insertDocuments(tx *gorm.DB, documents []models.Document) error {
	return tx.Create(&documents).Error
}

func (d *Database) handleTerms(tx *gorm.DB, uniqueTerms map[string]uint) error {
	termKeys := make([]string, 0, len(uniqueTerms))
	for term := range uniqueTerms {
		termKeys = append(termKeys, term)
	}

	var existingTerms []models.Term
	if err := tx.Where("text IN ?", termKeys).Find(&existingTerms).Error; err != nil {
		return err
	}

	existingTermMap := make(map[string]models.Term)
	for _, term := range existingTerms {
		existingTermMap[term.Text] = term
	}

	var newTerms []models.Term
	for term, docCount := range uniqueTerms {
		if _, exists := existingTermMap[term]; !exists {
			newTerms = append(newTerms, models.Term{Text: term, DocCount: docCount})
		}
	}
	if len(newTerms) > 0 {
		if err := tx.Create(&newTerms).Error; err != nil {
			return err
		}
	}

	for _, term := range existingTerms {
		increment := uniqueTerms[term.Text]
		if err := tx.Model(&term).Update("doc_count", gorm.Expr("doc_count + ?", increment)).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) insertDocumentTerms(tx *gorm.DB, documents []models.Document, termFreqs []map[string]uint) error {
	var allTerms []models.Term
	if err := tx.Find(&allTerms).Error; err != nil {
		return err
	}
	termTextToID := make(map[string]int)
	for _, term := range allTerms {
		termTextToID[term.Text] = term.ID
	}

	var documentTerms []models.DocumentTerm
	for i, doc := range documents {
		tf := termFreqs[i]
		for term, count := range tf {
			termID := termTextToID[term]
			documentTerms = append(documentTerms, models.DocumentTerm{
				DocumentID: doc.ID,
				TermID:     termID,
				Count:      count,
			})
		}
	}

	if len(documentTerms) > 0 {
		return tx.Create(&documentTerms).Error
	}
	return nil
}

func (d *Database) GetTotalDocuments() (int64, error) {
	var count int64
	err := d.db.Model(&models.Document{}).Count(&count).Error
	return count, err
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
