package repositories

import (
	"book-manager/internal/models/user"
	"book-manager/internal/utils/enums/error_codes"
	"book-manager/internal/utils/logger"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func (r *AuthRepository) dbCtx(requestId string) *gorm.DB {
	ctx := logger.ContextWithRequestID(context.Background(), requestId)
	return r.DB.WithContext(ctx)
}

func (r *AuthRepository) Create(u *user.User, requestId string) error {
	db := r.dbCtx(requestId)

	// Check if username already exists
	var existing user.User
	err := db.Where("LOWER(username) = LOWER(?)", u.Username).First(&existing).Error
	if err == nil {
		return error_codes.NewBookStoreError(error_codes.UserAlreadyExist, requestId)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return error_codes.ThrowException(err, requestId)
	}

	if err := db.Create(u).Error; err != nil {
		return error_codes.ThrowException(err, requestId)
	}
	return nil
}

func (r *AuthRepository) FindByUsername(username string, requestId string) (*user.User, error) {
	db := r.dbCtx(requestId)

	var u user.User
	err := db.Where("LOWER(username) = LOWER(?) AND status <> 'deleted'", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error_codes.NewBookStoreError(error_codes.InvalidCredential, requestId)
		}
		return nil, error_codes.ThrowException(err, requestId)
	}
	return &u, nil
}

func (r *AuthRepository) RevokeToken(token string, expiresAt interface{}, requestId string) error {
	db := r.dbCtx(requestId)

	revoked := user.RevokedToken{
		Token: token,
	}

	if t, ok := expiresAt.(time.Time); ok {
		revoked.ExpiresAt = t
	} else {
		revoked.ExpiresAt = time.Now().Add(7 * 24 * time.Hour) // fallback 7 days
	}

	return db.Create(&revoked).Error
}

func (r *AuthRepository) IsTokenRevoked(token string) bool {
	var count int64
	r.DB.Model(&user.RevokedToken{}).Where("token = ?", token).Count(&count)
	return count > 0
}
