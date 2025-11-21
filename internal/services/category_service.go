package services

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/models/common"
)

type ICategoryService interface {
	Create(req *common.Request[category.AddCategory]) common.Response[any]
}
