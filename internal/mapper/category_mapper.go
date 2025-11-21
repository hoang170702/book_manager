package mapper

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/models"
	"book-manager/internal/utils/enums"
	"time"
)

func CategoryMapper(add category.AddCategory) *models.Category {
	ctg := &models.Category{
		Name: add.Name,
	}
	ctg.CreatedBy = add.CreatedBy
	ctg.CreatedDate = time.Now()
	ctg.Status = enums.StatusActive
	return ctg
}
