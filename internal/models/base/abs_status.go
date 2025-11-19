package base

import (
	"book-manager/internal/utils/enums"

	"gorm.io/gorm"
)

type AbstractStatus struct {
	Status enums.Status `json:"status" gorm:"type:VARCHAR(20);default:'active'"`
}

func (b *AbstractStatus) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Status == "" {
		b.Status = enums.StatusActive
	}
	return
}

func (b *AbstractStatus) BeforeUpdate(tx *gorm.DB) (err error) {
	if !b.Status.IsValid() {
		b.Status = enums.StatusActive
	}
	return
}
