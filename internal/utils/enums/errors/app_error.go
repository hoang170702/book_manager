package errors

type AppError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id"`
	Err       error     `json:"-"`
}

func New(code ErrorCode, cause error, requestId string) *AppError {
	return &AppError{
		Code:      code,
		Message:   getMessage(code),
		RequestID: requestId,
		Err:       cause,
	}
}

func (e *AppError) Error() string {
	return string(e.Code) + ": " + e.Message
}

func ThrowBookStoreException(code ErrorCode, cause error, requestId string) *AppError {
	return New(code, cause, requestId)
}
