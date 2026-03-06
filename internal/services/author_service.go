package services

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
)

type IAuthorService interface {
	Create(req *common.Request[author.AddAuthor]) common.Response[any]
	GetOne(req *common.Request[author.GetOneAuthor]) common.Response[author.AuthorResponse]
	GetAll(req *common.Request[any]) common.Response[[]author.AuthorResponse]
	Update(req *common.Request[author.UpdateAuthor]) common.Response[any]
	Delete(req *common.Request[author.DeleteAuthor]) common.Response[any]
}
