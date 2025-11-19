package models

import (
	"book-manager/internal/models/base"
)

type Author struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" gorm:"default:active"`
	base.AbsTimestamp
}
