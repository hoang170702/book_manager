package models

import (
	"book-manager/internal/models/base"
)

type Category struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" binding:"required"`
	base.AbstractStatus
	base.AbsTimestamp
}
