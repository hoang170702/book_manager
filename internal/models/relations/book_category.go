package relations

import (
	"book-manager/internal/models"
	"book-manager/internal/models/base"
	"book-manager/internal/models/book"
)

type BookCategory struct {
	Id         int `gorm:"primaryKey;autoIncrement" json:"id"`
	BookID     int `json:"book_id"`
	CategoryID int `json:"category_id"`

	Book     book.Book       `gorm:"foreignKey:BookID;"`
	Category models.Category `gorm:"foreignKey:CategoryID;"`

	base.AbstractStatus
	base.AbsTimestamp
}
