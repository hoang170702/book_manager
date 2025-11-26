package services

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
)

type IAuthorService interface {
	Create(req *common.Request[author.AddAuthor]) common.Response[any]
	GetOne(req *common.Request[author.GetOneAuthor]) common.Response[models.Author]
	GetAll(req *common.Request[any]) common.Response[[]models.Author]
	Update(req *common.Request[author.UpdateAuthor]) common.Response[any]
	Delete(req *common.Request[author.DeleteAuthor]) common.Response[any]
}
