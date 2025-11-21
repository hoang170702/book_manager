package impl

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/mapper"
	"book-manager/internal/models/common"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) services.ICategoryService {
	return &CategoryService{Repo: repo}
}

func (c CategoryService) Create(req *common.Request[category.AddCategory]) common.Response[any] {
	categoryReq := mapper.CategoryMapper(req.Data)

	CategoryReq := mapper.RequestMapper(req, categoryReq)

	ok, err := c.Repo.Create(CategoryReq)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[any](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		})
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest)
	}

	return utils.BuildResponse[any](nil, error_codes.Success)
}
