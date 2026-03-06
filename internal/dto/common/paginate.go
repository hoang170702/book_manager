package common

// Paginate is used in Request DTOs for pagination
type Paginate struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
