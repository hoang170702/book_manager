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
	Repo *repositories.AuthorRepository
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
