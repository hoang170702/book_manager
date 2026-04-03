package user

import "book-manager/internal/models/base"

// User represents the user entity for authentication
type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"type:VARCHAR(100);uniqueIndex;not null"`
	Password string `json:"-" gorm:"type:VARCHAR(255);not null"`
	base.AbstractStatus
	base.AbsTimestamp
}
