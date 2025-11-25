package services

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
)

type IAuthorService interface {
	Create(req *common.Request[author.AddAuthor]) common.Response[any]
}
