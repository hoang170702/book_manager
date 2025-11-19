package models

type Author struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" binding:"required"`
}
