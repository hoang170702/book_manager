package base

import "time"

type AbsTimestamp struct {
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedDate time.Time `json:"created_date" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedDate time.Time `json:"updated_date"`
}
