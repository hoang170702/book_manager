package repositories

import (
	"book-manager/internal/dto/author"
	"book-manager/internal/dto/common"
	"book-manager/internal/models"
)

type IAuthorRepository interface {
	Create(request common.Request[*models.Author]) (bool, error)
	GetOne(request *common.Request[author.GetOneAuthor]) (models.Author, error)
	GetAll(request *common.Request[any]) ([]models.Author, error)
	Update(request *common.Request[author.UpdateAuthor], user string) (bool, error)
	Delete(request *common.Request[author.DeleteAuthor], user string) (bool, error)
}
