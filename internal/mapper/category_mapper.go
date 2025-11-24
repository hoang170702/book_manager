package mapper

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums"
	"time"
)

func CategoryMapper(add category.AddCategory, user string) *models.Category {
	ctg := &models.Category{
		Name: add.Name,
	}
	ctg.CreatedBy = user
	ctg.CreatedDate = time.Now()
	ctg.Status = enums.StatusActive
	return ctg
}
