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
	Repo repositories.ICategoryRepository
}

func (c CategoryService) Delete(req *common.Request[category.DeleteCategory]) common.Response[any] {
	ok, err := c.Repo.Delete(req, "Anonymous")

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[any](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		}, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func (c CategoryService) Update(req *common.Request[category.UpdateCategory]) common.Response[any] {
	ok, err := c.Repo.Update(req, "Anonymous")

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[any](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		}, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func (c CategoryService) GetAll(req *common.Request[any]) common.Response[[]models.Category] {
	data, err := c.Repo.GetAll(req)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[[]models.Category](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		}, req.RequestId)
	}
	return utils.BuildResponse[[]models.Category](data, error_codes.Success, req.RequestId)
}

func (c CategoryService) GetOne(req *common.Request[category.GetOneCategory]) common.Response[models.Category] {

	data, err := c.Repo.GetOne(req)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[models.Category](models.Category{}, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		}, req.RequestId)
	}

	return utils.BuildResponse[models.Category](data, error_codes.Success, req.RequestId)
}

func (c CategoryService) Create(req *common.Request[category.AddCategory]) common.Response[any] {
	categoryReq := mapper.CategoryMapper(req.Data, "Anonymous")

	CategoryReq := mapper.RequestMapper(req, categoryReq)

	ok, err := c.Repo.Create(CategoryReq)

	if err != nil {
		appErr := err.(*error_codes.AppError)
		return utils.BuildResponse[any](nil, error_codes.ErrorCode{
			Code: appErr.Code,
			Msg:  appErr.Message,
		}, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func NewCategoryService(repo repositories.ICategoryRepository) services.ICategoryService {
	return &CategoryService{Repo: repo}
}
