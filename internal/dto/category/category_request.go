package category

import "book-manager/internal/dto"

type AddCategory struct {
	Name string `json:"name" validate:"required"`
	dto.BaseCreateDto
}
