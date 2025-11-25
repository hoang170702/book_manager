package repositories

import (
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums/error_codes"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	DB *gorm.DB
}

func (r *AuthorRepository) Add(request common.Request[*models.Author]) (bool, error) {

	author := request.Data
	var existAuthor models.Author

	result := r.DB.Where("name = ?", author.Name).First(&existAuthor)

	if result.RowsAffected > 0 {
		return false, error_codes.NewBookStoreError(error_codes.AuthorAlreadyExist, request.RequestId)
	}

	if err := r.DB.Create(author).Error; err != nil {
		return false, error_codes.ThrowException(err, request.RequestId)
	}

	return true, nil
}
