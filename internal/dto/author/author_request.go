package author

type AddAuthor struct {
	Name string `json:"name" validate:"required"`
}

type GetOneAuthor struct {
	Id int `json:"id" validate:"required"`
}

type UpdateAuthor struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteAuthor struct {
	Id int `json:"id" validate:"required"`
}
