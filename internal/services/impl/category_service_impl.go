package impl

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/mapper"
	"book-manager/internal/models"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

func (c CategoryService) GetAll(req *common.Request[any]) common.Response[[]models.Category] {
	data, err := c.Repo.GetAll(req)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[[]models.Category](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		})
	}
	return utils.BuildResponse[[]models.Category](data, error_codes.Success)
}

func (c CategoryService) GetOne(req *common.Request[category.GetOneCategory]) common.Response[models.Category] {

	data, err := c.Repo.GetOne(req)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[models.Category](models.Category{}, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		})
	}

	return utils.BuildResponse[models.Category](data, error_codes.Success)
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

func NewCategoryService(repo *repositories.CategoryRepository) services.ICategoryService {
	return &CategoryService{Repo: repo}
}
