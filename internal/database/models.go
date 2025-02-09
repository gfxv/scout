package main

type Document struct {
	ID    int
	Path  string  `gorm:"unique"`
	Terms []*Term `gorm:"many2many:document_indexes;"`
}

type Term struct {
	ID   int
	Term string `gorm:"unique"`
	Freq float32
}

type DocumentIndex struct {
	TermID     int `gorm:"primaryKey"`
	DocumentID int `gorm:"primaryKey"`
	Freq       float32
}
