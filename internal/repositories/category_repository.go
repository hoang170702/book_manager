package repositories

import (
	"book-manager/internal/models"
	"book-manager/internal/models/common"
	"book-manager/internal/utils/enums/error_codes"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) Create(request common.Request[*models.Category]) (bool, error) {
	var category = request.Data
	var existCategory models.Category
	err := r.DB.Where("name = ?", category.Name).First(&existCategory).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.DB.Create(category).Error; err != nil {
			return false, error_codes.ThrowException(err, request.RequestId)
		}
		return true, nil
	}

	if err != nil {
		return false, error_codes.ThrowException(err, request.RequestId)
	}

	return false, error_codes.ThrowBookStoreException(error_codes.CreateCategoryFailed, request.RequestId)
}
