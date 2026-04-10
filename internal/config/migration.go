package config

import (
	"book-manager/internal/models"
	"book-manager/internal/models/book"
	"book-manager/internal/models/relations"
	"book-manager/internal/models/user"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	modelsList := []interface{}{
		&models.Author{},
		&models.Category{},
		&book.Book{},
		&relations.BookCategory{},
		&user.User{},
		&user.RevokedToken{},
	}

	for _, m := range modelsList {
		if err := db.AutoMigrate(m); err != nil {
			log.Fatalf("❌ Migration failed for %T: %v", m, err)
		}

		if !db.Migrator().HasTable(m) {
			log.Fatalf("❌ Table does not exist after migration: %T", m)
		}

		log.Printf("✔ Migration OK: %T", m)
	}
}
