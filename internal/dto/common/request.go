package common

type Request[T any] struct {
	RequestId   string    `json:"request_id"`
	RequestTime string    `json:"request_time"`
	Data        T         `json:"data"`
	Paginate    *Paginate `json:"paginate"`
}
