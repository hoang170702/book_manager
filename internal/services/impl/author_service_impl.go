package impl

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
)

type AuthorService struct {
	Repo *repositories.AuthorRepository
}

func (a AuthorService) Create(req *common.Request[author.AddAuthor]) common.Response[any] {
	//TODO implement me
	panic("implement me")
}

func NewAuthorService(repo *repositories.AuthorRepository) services.IAuthorService {
	return &AuthorService{Repo: repo}
}
