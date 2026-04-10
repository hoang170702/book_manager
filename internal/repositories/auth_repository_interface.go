package repositories

import (
	"book-manager/internal/models/user"
)

type IAuthRepository interface {
	Create(u *user.User, requestId string) error
	FindByUsername(username string, requestId string) (*user.User, error)
	RevokeToken(token string, expiresAt interface{}, requestId string) error
	IsTokenRevoked(token string) bool
}
