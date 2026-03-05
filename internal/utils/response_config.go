package utils

import (
	"book-manager/internal/dto/common"
	"book-manager/internal/utils/enums/error_codes"
	"book-manager/internal/utils/logger"
	"time"
)

func BuildResponse[T any](data T, code error_codes.ErrorCode, requestId string) common.Response[T] {

	if code.Code != error_codes.Success.Code {
		logger.Warn("RESPONSE", nil, "[RequestId: %s] code=%s, msg=%s", requestId, code.Code, code.Msg)
	}

	return common.Response[T]{
		ResponseId:   requestId,
		ResponseCode: code.Code,
		ResponseMsg:  code.Msg,
		Data:         data,
		ResponseTime: time.Now().Format(time.RFC3339),
	}
}
