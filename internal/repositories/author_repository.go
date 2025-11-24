package repositories

import "gorm.io/gorm"

type AuthorRepository struct {
	DB *gorm.DB
}
