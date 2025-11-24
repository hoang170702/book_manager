package common

import "book-manager/internal/models"

type Request[T any] struct {
	RequestId   string           `json:"request_id"`
	RequestTime string           `json:"request_time"`
	Data        T                `json:"data"`
	Paginate    *models.Paginate `json:"paginate"`
}
