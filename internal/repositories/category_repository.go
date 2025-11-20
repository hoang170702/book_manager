package repositories

import (
	"book-manager/internal/models"
	"book-manager/internal/utils/enums/error_codes"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) Create(category *models.Category, requestId string) error {

	if err := r.DB.Create(category).Error; err != nil {
		return error_codes.ThrowException(err, requestId)
	}
	return nil
}

func (r *CategoryRepository) FindByName(name string, requestId string) (bool, error) {
	var category models.Category
	err := r.DB.Where("name = ?", name).First(&category).Error

	if err != nil {
		return false, error_codes.ThrowException(err, requestId)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}

	return true, nil
}
