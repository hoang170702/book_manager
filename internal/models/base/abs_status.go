package base

import (
	"book-manager/internal/utils/enums"
	"errors"

	"gorm.io/gorm"
)

type AbstractStatus struct {
	Status enums.Status `json:"status" gorm:"type:VARCHAR(20);default:'active'"`
}

func (b *AbstractStatus) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Status == "" {
		b.Status = enums.StatusActive
	}
	if !b.Status.IsValid() {
		return errors.New("invalid status value")
	}
	return
}

func (b *AbstractStatus) BeforeUpdate(tx *gorm.DB) (err error) {
	if !b.Status.IsValid() {
		return errors.New("invalid status value")
	}
	return
}
