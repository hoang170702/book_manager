package repositories

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
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

func (r *CategoryRepository) GetOne(request *common.Request[category.GetOneCategory]) (models.Category, error) {
	var category = request.Data
	var existCategory models.Category
	err := r.DB.Where("id = ?", category.Id).First(&existCategory).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, error_codes.ThrowBookStoreException(error_codes.CategoryNotFound, request.RequestId)
		}
		return models.Category{}, error_codes.ThrowException(err, request.RequestId)
	}

	return existCategory, nil
}

func (r *CategoryRepository) GetAll(request *common.Request[any]) ([]models.Category, error) {
	var categories []models.Category

	err := r.DB.Find(&categories).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_codes.ThrowBookStoreException(error_codes.CategoryNotFound, request.RequestId)
		}
		return nil, error_codes.ThrowException(err, request.RequestId)
	}
	return categories, nil
}
