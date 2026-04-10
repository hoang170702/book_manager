package user

import "time"

// RevokedToken stores blacklisted JWT tokens (after logout)
type RevokedToken struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Token     string    `json:"token" gorm:"type:TEXT;not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"` // auto-cleanup after expiry
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
