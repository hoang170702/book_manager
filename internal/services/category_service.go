package services

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
)

type ICategoryService interface {
	Create(req *common.Request[category.AddCategory]) common.Response[any]
	GetOne(req *common.Request[category.GetOneCategory]) common.Response[category.CategoryResponse]
	GetAll(req *common.Request[any]) common.Response[[]category.CategoryResponse]
	Update(req *common.Request[category.UpdateCategory]) common.Response[any]
	Delete(req *common.Request[category.DeleteCategory]) common.Response[any]
}
