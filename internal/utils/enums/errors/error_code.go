package errors

type ErrorCode string

const (
	CodeApproved   ErrorCode = "00"
	CodeBadRequest ErrorCode = "99"

	CategoryNameExists ErrorCode = "01"
)
