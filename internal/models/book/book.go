package book

import (
	"book-manager/internal/models"
	"book-manager/internal/models/base"
)

type Book struct {
	Id    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title" binding:"required"`
	Price int    `json:"price" binding:"required"`

	// ---- Author (1-N) ----
	AuthorID int           `json:"author_id" binding:"required"`
	Author   models.Author `json:"author" gorm:"foreignKey:AuthorID"`

	// ---- Category (N-N) ----
	Categories []models.Category `gorm:"many2many:book_categories;"`

	Year string `json:"year" binding:"required"`
	base.AbstractStatus
	base.AbsTimestamp
}
