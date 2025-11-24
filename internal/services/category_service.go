package services

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
)

type ICategoryService interface {
	Create(req *common.Request[category.AddCategory]) common.Response[any]
	GetOne(req *common.Request[category.GetOneCategory]) common.Response[models.Category]
	GetAll(req *common.Request[any]) common.Response[[]models.Category]
	Update(req *common.Request[category.UpdateCategory]) common.Response[any]
	Delete(req *common.Request[category.DeleteCategory]) common.Response[any]
}
