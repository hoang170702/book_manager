package relations

import (
	"book-manager/internal/models/base"
)

type BookCategory struct {
	BookID     int    `gorm:"primaryKey" json:"book_id"`
	CategoryID int    `gorm:"primaryKey" json:"category_id"`
	Status     string `json:"status" gorm:"default:active"`
	base.AbsTimestamp
}
