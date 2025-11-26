package impl

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/mapper"
	"book-manager/internal/models"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/utils"
	"book-manager/internal/utils/enums/error_codes"
	"errors"
)

type AuthorService struct {
	Repo *repositories.AuthorRepository
}

func (a AuthorService) Delete(req *common.Request[author.DeleteAuthor]) common.Response[any] {
	err := a.Repo.Delete(req, "Anonymous")
	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			})
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest)
	}
	return utils.BuildResponse[any](nil, error_codes.Success)
}

func (a AuthorService) Update(req *common.Request[author.UpdateAuthor]) common.Response[any] {
	err := a.Repo.Update(req, "Anonymous")

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			})
		}
		return utils.BuildResponse[any](nil, error_codes.BadRequest)
	}
	return utils.BuildResponse[any](nil, error_codes.Success)
}

func (a AuthorService) GetAll(req *common.Request[any]) common.Response[[]models.Author] {
	data, err := a.Repo.GetAll(req)
	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[[]models.Author]([]models.Author{}, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			})
		}
		return utils.BuildResponse[[]models.Author]([]models.Author{}, error_codes.BadRequest)
	}
	return utils.BuildResponse[[]models.Author](data, error_codes.Success)
}

func (a AuthorService) GetOne(req *common.Request[author.GetOneAuthor]) common.Response[models.Author] {
	data, err := a.Repo.GetOne(req)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[models.Author](models.Author{}, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			})
		}
		return utils.BuildResponse[models.Author](models.Author{}, error_codes.BadRequest)
	}
	return utils.BuildResponse[models.Author](data, error_codes.Success)
}

func (a AuthorService) Create(req *common.Request[author.AddAuthor]) common.Response[any] {
	authorEntity := mapper.AuthorMapper(req.Data, "Anonymous")
	request := mapper.RequestMapper(req, authorEntity)

	err := a.Repo.Create(request)

	if err != nil {
		var appErr *error_codes.AppError
		if errors.As(err, &appErr) {
			return utils.BuildResponse[any](nil, error_codes.ErrorCode{
				Code: appErr.Code,
				Msg:  appErr.Message,
			})
		}

		return utils.BuildResponse[any](nil, error_codes.BadRequest)
	}
	return utils.BuildResponse[any](nil, error_codes.Success)

}

func NewAuthorService(repo *repositories.AuthorRepository) services.IAuthorService {
	return &AuthorService{Repo: repo}
}
