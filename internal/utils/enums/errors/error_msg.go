package errors

var errorMsg = map[ErrorCode]string{
	CodeApproved:   "Approved",
	CodeBadRequest: "Bad Request",

	CategoryNameExists: "Category name already exists",
}

func getMessage(code ErrorCode) string {
	if msg, ok := errorMsg[code]; ok {
		return msg
	}
	return "Unknown error"
}
