package repositories

import (
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
