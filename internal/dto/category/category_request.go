package category

type AddCategory struct {
	Name string `json:"name" validate:"required"`
}

type GetOneCategory struct {
	Id int `json:"id" validate:"required"`
}

type UpdateCategory struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteCategory struct {
	Id int `json:"id" validate:"required"`
}
