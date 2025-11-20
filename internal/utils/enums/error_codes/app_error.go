package error_codes

import "fmt"

type AppError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Err       error  `json:"-"`
}

func NewBookStoreError(code ErrorCode, requestId string) *AppError {
	return &AppError{
		Code:      code.Code,
		Message:   code.Msg,
		RequestID: requestId,
	}
}

func NewExceptionError(code ErrorCode, cause error, requestId string) *AppError {
	return &AppError{
		Code:      code.Code,
		Message:   code.Msg,
		RequestID: requestId,
		Err:       cause,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func ThrowBookStoreException(code ErrorCode, requestId string) *AppError {
	return NewBookStoreError(code, requestId)
}

func ThrowException(cause error, requestId string) *AppError {
	return NewExceptionError(BadRequest, cause, requestId)
}
