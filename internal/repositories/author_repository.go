package repositories

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums/error_codes"
	"errors"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	DB *gorm.DB
}

func (r *AuthorRepository) Create(request common.Request[*models.Author]) error {

	author := request.Data
	var existAuthor models.Author

	err := r.DB.Where("LOWER(name) = LOWER(?)", author.Name).First(&existAuthor).Error

	if err == nil {
		return error_codes.NewBookStoreError(error_codes.AuthorAlreadyExist, request.RequestId)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_codes.ThrowException(err, request.RequestId)
	}

	if err := r.DB.Create(author).Error; err != nil {
		return error_codes.ThrowException(err, request.RequestId)
	}
	return nil
}

func (r *AuthorRepository) GetOne(request *common.Request[author.GetOneAuthor]) (models.Author, error) {
	author := request.Data
	var existAuthor models.Author

	err := r.DB.Where("id = ?", author.Id).First(&existAuthor).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Author{}, error_codes.NewBookStoreError(error_codes.AuthorNotFound, request.RequestId)
		}
		return models.Author{}, error_codes.ThrowException(err, request.RequestId)
	}
	return existAuthor, nil
}

func (r *AuthorRepository) GetAll(request *common.Request[any]) ([]models.Author, error) {
	var authors []models.Author

	err := r.DB.Model(&models.Author{}).Find(&authors).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Author{}, error_codes.NewBookStoreError(error_codes.AuthorNotFound, request.RequestId)
		}
		return []models.Author{}, error_codes.ThrowException(err, request.RequestId)
	}

	return authors, nil
}
