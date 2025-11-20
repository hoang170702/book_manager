package repositories

import (
	"book-manager/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.DB.Create(category).Error
}
