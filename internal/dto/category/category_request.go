package category

import "book-manager/internal/dto"

type AddCategory struct {
	Name string `json:"name" validate:"required"`
	dto.BaseCreateDto
}

type GetOneCategory struct {
	Id int `json:"id" validate:"required"`
}
