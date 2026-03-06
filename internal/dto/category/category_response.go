package category

// CategoryResponse is the API response DTO — only exposes fields client needs
type CategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
