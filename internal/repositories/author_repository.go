package repositories

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums"
	"book-manager/internal/utils/enums/error_codes"
	"book-manager/internal/utils/logger"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	DB *gorm.DB
}

// dbCtx creates a GORM session with request ID in context for tracing
func (r *AuthorRepository) dbCtx(requestId string) *gorm.DB {
	ctx := logger.ContextWithRequestID(context.Background(), requestId)
	return r.DB.WithContext(ctx)
}

func (r *AuthorRepository) Create(request common.Request[*models.Author]) error {
	db := r.dbCtx(request.RequestId)

	author := request.Data
	var existAuthor models.Author

	err := db.Where("LOWER(name) = LOWER(?)", author.Name).First(&existAuthor).Error

	if err == nil {
		return error_codes.NewBookStoreError(error_codes.AuthorAlreadyExist, request.RequestId)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_codes.ThrowException(err, request.RequestId)
	}

	if err := db.Create(author).Error; err != nil {
		return error_codes.ThrowException(err, request.RequestId)
	}
	return nil
}

func (r *AuthorRepository) GetOne(request *common.Request[author.GetOneAuthor]) (models.Author, error) {
	db := r.dbCtx(request.RequestId)
	author := request.Data
	var existAuthor models.Author

	err := db.Where("id = ? and status <> ?", author.Id, enums.StatusDeleted).First(&existAuthor).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Author{}, error_codes.NewBookStoreError(error_codes.AuthorNotFound, request.RequestId)
		}
		return models.Author{}, error_codes.ThrowException(err, request.RequestId)
	}
	return existAuthor, nil
}

func (r *AuthorRepository) GetAll(request *common.Request[any]) ([]models.Author, error) {
	db := r.dbCtx(request.RequestId)
	var authors []models.Author

	err := db.Where("status != ?", enums.StatusDeleted).
		Find(&authors).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Author{}, error_codes.NewBookStoreError(error_codes.AuthorNotFound, request.RequestId)
		}
		return []models.Author{}, error_codes.ThrowException(err, request.RequestId)
	}

	return authors, nil
}

func (r *AuthorRepository) Update(request *common.Request[author.UpdateAuthor], user string) error {
	db := r.dbCtx(request.RequestId)
	data := request.Data
	var existAuthor models.Author

	err := db.Model(&models.Author{}).
		Where("name = ? and  id <> ? and status <> ?", data.Name, data.Id, enums.StatusDeleted).
		First(&existAuthor).Error

	if err == nil {
		return error_codes.NewBookStoreError(error_codes.AuthorAlreadyExist, request.RequestId)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_codes.NewExceptionError(error_codes.BadRequest, err, request.RequestId)
	}

	update := map[string]interface{}{
		"name":         data.Name,
		"updated_by":   user,
		"updated_date": time.Now(),
	}

	err = db.Model(&models.Author{}).
		Where("id = ?", data.Id).
		Updates(update).Error

	if err != nil {
		return error_codes.NewExceptionError(error_codes.BadRequest, err, request.RequestId)
	}

	return nil
}

func (r *AuthorRepository) Delete(request *common.Request[author.DeleteAuthor], user string) error {
	db := r.dbCtx(request.RequestId)
	data := request.Data

	deleted := map[string]interface{}{
		"status":       enums.StatusDeleted,
		"updated_by":   user,
		"updated_date": time.Now(),
	}

	err := db.Model(&models.Author{}).
		Where("id = ? and status <> ?", data.Id, enums.StatusDeleted).
		UpdateColumns(deleted).Error

	if err != nil {
		return error_codes.NewExceptionError(error_codes.BadRequest, err, request.RequestId)
	}

	return nil
}
