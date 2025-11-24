package mapper

import (
	"book-manager/internal/dto/common"
)

func RequestMapper[T any, R any](req *common.Request[T], data R) common.Request[R] {
	return common.Request[R]{
		RequestId:   req.RequestId,
		RequestTime: req.RequestTime,
		Paginate:    req.Paginate,
		Data:        data,
	}
}
