package repositories

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
)

type ICategoryRepository interface {
	Create(request common.Request[*models.Category]) (bool, error)
	GetOne(request *common.Request[category.GetOneCategory]) (models.Category, error)
	GetAll(request *common.Request[any]) ([]models.Category, error)
	Update(request *common.Request[category.UpdateCategory], user string) (bool, error)
	Delete(request *common.Request[category.DeleteCategory], user string) (bool, error)
}
