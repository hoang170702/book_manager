package impl

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/mapper"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
	"errors"
)

type AuthorService struct {
	Repo repositories.IAuthorRepository
}

func (a AuthorService) Delete(req *common.Request[author.DeleteAuthor], user string) common.Response[any] {
	ok, err := a.Repo.Delete(req, user)

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

func (a AuthorService) Update(req *common.Request[author.UpdateAuthor], user string) common.Response[any] {
	ok, err := a.Repo.Update(req, user)

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

func (a AuthorService) GetAll(req *common.Request[any]) common.Response[[]author.AuthorResponse] {
	data, err := a.Repo.GetAll(req)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[[]author.AuthorResponse](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[[]author.AuthorResponse](nil, error_codes.BadRequest, req.RequestId)
	}

	// Map entities to response DTOs
	var result []author.AuthorResponse
	for _, item := range data {
		result = append(result, author.AuthorResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return utils.BuildResponse[[]author.AuthorResponse](result, error_codes.Success, req.RequestId)
}

func (a AuthorService) GetOne(req *common.Request[author.GetOneAuthor]) common.Response[author.AuthorResponse] {

	data, err := a.Repo.GetOne(req)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[author.AuthorResponse](author.AuthorResponse{}, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			}, req.RequestId)
		}
		return utils.BuildResponse[author.AuthorResponse](author.AuthorResponse{}, error_codes.BadRequest, req.RequestId)
	}

	// Map entity to response DTO
	result := author.AuthorResponse{
		Id:   data.Id,
		Name: data.Name,
	}

	return utils.BuildResponse[author.AuthorResponse](result, error_codes.Success, req.RequestId)
}

func (a AuthorService) Create(req *common.Request[author.AddAuthor], user string) common.Response[any] {
	authorEntity := mapper.AuthorMapper(req.Data, user)

	request := mapper.RequestMapper(req, authorEntity)

	ok, err := a.Repo.Create(request)

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

func NewAuthorService(repo repositories.IAuthorRepository) services.IAuthorService {
	return &AuthorService{Repo: repo}
}
