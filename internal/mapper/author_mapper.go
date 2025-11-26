package mapper

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums"
	"time"
)

func AuthorMapper(add author.AddAuthor, user string) *models.Author {
	ato := &models.Author{Name: add.Name}
	ato.CreatedBy = user
	ato.CreatedDate = time.Now()
	ato.Status = enums.StatusActive
	return ato
}
