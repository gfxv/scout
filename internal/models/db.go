package models

type Document struct {
	ID         int    `gorm:"primaryKey"`
	Path       string `gorm:"unique;index"`
	TotalTerms uint
	Terms      []*Term `gorm:"many2many:document_terms;"`
}

type Term struct {
	ID       int    `gorm:"primaryKey"`
	Text     string `gorm:"unique;index;collate:NOCASE"`
	DocCount uint
}

type DocumentTerm struct {
	DocumentID int `gorm:"primaryKey"`
	TermID     int `gorm:"primaryKey"`
	Count      uint
}
