package relations

import (
	"book-manager/internal/models/base"
)

type BookCategory struct {
	BookID     int `gorm:"primaryKey" json:"book_id"`
	CategoryID int `gorm:"primaryKey" json:"category_id"`
	base.AbstractStatus
	base.AbsTimestamp
}
