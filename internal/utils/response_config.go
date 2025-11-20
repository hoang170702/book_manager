package utils

import (
	"book-manager/internal/models/common"
	"book-manager/internal/utils/enums/error_codes"
	"time"
)

func BuildResponse[T any](data T, code error_codes.ErrorCode) common.Response[T] {
	return common.Response[T]{
		ResponseCode: code.Code,
		ResponseMsg:  code.Msg,
		Data:         data,
		ResponseTime: time.Now().Format(time.RFC3339),
	}
}
