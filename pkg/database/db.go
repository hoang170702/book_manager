package database

import (
	"book-manager/internal/models"
	"book-manager/internal/models/book"
	"book-manager/internal/models/relations"
	"book-manager/internal/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		utils.GetEnv("DB_HOST", "localhost"),
		utils.GetEnv("DB_USER", "postgres"),
		utils.GetEnv("DB_PASSWORD", "postgres"),
		utils.GetEnv("DB_NAME", "postgres"),
		utils.GetEnv("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Không thể kết nối đến PostgreSQL: " + err.Error())
	}
	log.Println("Kết nối DB thành công!")

	migrate()
}

func migrate() {
	modelsList := []interface{}{
		&models.Author{},
		&models.Category{},
		&book.Book{},
		&relations.BookCategory{},
	}

	for _, m := range modelsList {
		if err := DB.AutoMigrate(m); err != nil {
			log.Fatalf("❌ Migration thất bại ở %T: %v", m, err)
		}

		if !DB.Migrator().HasTable(m) {
			log.Fatalf("❌ Bảng không tồn tại sau migrate: %T", m)
		}

		log.Printf("✔ Migration OK: %T", m)
	}
}
