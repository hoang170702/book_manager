package repositories

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums"
	"book-manager/internal/utils/enums/error_codes"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func (r *CategoryRepository) Create(request common.Request[*models.Category]) (bool, error) {
	var category = request.Data
	var existCategory models.Category
	result := r.DB.Where("name = ?", category.Name).First(&existCategory)

	if result.RowsAffected > 0 {
		return false, error_codes.NewBookStoreError(error_codes.CategoryAlreadyExist, request.RequestId)
	}

	if err := r.DB.Create(category).Error; err != nil {
		return false, error_codes.ThrowException(err, request.RequestId)
	}

	return true, nil
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

func (r *CategoryRepository) Update(request *common.Request[category.UpdateCategory], user string) (bool, error) {
	data := request.Data
	var existCategory models.Category

	check := r.DB.Model(&models.Category{}).
		Where("name = ? AND id <> ?", data.Name, data.Id).
		First(&existCategory)

	if check.RowsAffected > 0 {
		return false, error_codes.NewBookStoreError(error_codes.CategoryAlreadyExist, request.RequestId)
	}

	update := map[string]interface{}{
		"name":         data.Name,
		"updated_by":   user,
		"updated_date": time.Now(),
	}

	result := r.DB.Model(&models.Category{}).
		Where("id = ?", data.Id).
		Updates(update)

	if result.Error != nil {
		return false, error_codes.ThrowException(result.Error, request.RequestId)
	}

	if result.RowsAffected == 0 {
		return false, error_codes.ThrowBookStoreException(error_codes.CategoryNotFound, request.RequestId)
	}

	return true, nil
}

func (r *CategoryRepository) Delete(request *common.Request[category.DeleteCategory], user string) (bool, error) {
	data := request.Data

	deleted := map[string]interface{}{
		"status":       enums.StatusDeleted,
		"updated_by":   user,
		"updated_date": time.Now(),
	}

	result := r.DB.Model(&models.Category{}).
		Where("id = ?", data.Id).
		UpdateColumns(deleted)

	if result.Error != nil {
		return false, error_codes.ThrowException(result.Error, request.RequestId)
	}

	if result.RowsAffected == 0 {
		return false, error_codes.ThrowBookStoreException(error_codes.CategoryNotFound, request.RequestId)
	}
	return true, nil
}
