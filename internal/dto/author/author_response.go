package author

// AuthorResponse is the API response DTO — only exposes fields client needs
type AuthorResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
