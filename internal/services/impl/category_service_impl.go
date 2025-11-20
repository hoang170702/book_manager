package impl

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/models"
	"book-manager/internal/models/common"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils/enums/error_codes"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) services.ICategoryService {
	return &CategoryService{Repo: repo}
}

func (c CategoryService) Create(req *common.Request[category.AddCategory]) error {
	categoryReq := req.Data

	if c.Repo.FindByName(categoryReq.Name) == false {
		return error_codes.ThrowBookStoreException(error_codes.CategoryAlreadyExist, req.RequestId)
	}
	categoryM := models.Category{
		Name:        categoryReq.Name,
		CreatedBy:   categoryReq.CreatedBy,
		CreatedDate: categoryReq.CreatedDate,
		Status:      categoryReq.Status,
	}
	return c.Repo.Create(categoryM, req.RequestId)
}
