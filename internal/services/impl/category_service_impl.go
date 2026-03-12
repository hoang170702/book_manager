package impl

import (
	"book-manager/internal/dto/category"
	"book-manager/internal/dto/common"
	"book-manager/internal/mapper"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
	"errors"
)

type CategoryService struct {
	Repo repositories.ICategoryRepository
}

func (c CategoryService) Delete(req *common.Request[category.DeleteCategory], user string) common.Response[any] {
	ok, err := c.Repo.Delete(req, user)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func (c CategoryService) Update(req *common.Request[category.UpdateCategory], user string) common.Response[any] {
	ok, err := c.Repo.Update(req, user)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func (c CategoryService) GetAll(req *common.Request[any]) common.Response[[]category.CategoryResponse] {
	data, err := c.Repo.GetAll(req)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[[]category.CategoryResponse](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[[]category.CategoryResponse](nil, error_codes.BadRequest, req.RequestId)
	}

	// Map entities to response DTOs
	var result []category.CategoryResponse
	for _, item := range data {
		result = append(result, category.CategoryResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return utils.BuildResponse[[]category.CategoryResponse](result, error_codes.Success, req.RequestId)
}

func (c CategoryService) GetOne(req *common.Request[category.GetOneCategory]) common.Response[category.CategoryResponse] {

	data, err := c.Repo.GetOne(req)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[category.CategoryResponse](category.CategoryResponse{}, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[category.CategoryResponse](category.CategoryResponse{}, error_codes.BadRequest, req.RequestId)
	}

	// Map entity to response DTO
	result := category.CategoryResponse{
		Id:   data.Id,
		Name: data.Name,
	}

	return utils.BuildResponse[category.CategoryResponse](result, error_codes.Success, req.RequestId)
}

func (c CategoryService) Create(req *common.Request[category.AddCategory], user string) common.Response[any] {
	categoryReq := mapper.CategoryMapper(req.Data, user)

	CategoryReq := mapper.RequestMapper(req, categoryReq)

	ok, err := c.Repo.Create(CategoryReq)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	if !ok {
		return utils.BuildResponse[any](nil, error_codes.BadRequest, req.RequestId)
	}

	return utils.BuildResponse[any](nil, error_codes.Success, req.RequestId)
}

func NewCategoryService(repo repositories.ICategoryRepository) services.ICategoryService {
	return &CategoryService{Repo: repo}
}
