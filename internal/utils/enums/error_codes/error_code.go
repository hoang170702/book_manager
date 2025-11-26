package error_codes

type ErrorCode struct {
	Code string
	Msg  string
}

var (
	Success    = ErrorCode{"00", "Success"}
	BadRequest = ErrorCode{"400", "Bad Request"}

	CategoryAlreadyExist = ErrorCode{"01", "Category already exist"}
	CategoryNotFound     = ErrorCode{"02", "Category not found"}

	AuthorAlreadyExist   = ErrorCode{"03", "Author already exist"}
	AuthorNotFound       = ErrorCode{"04", "Author not found"}
	AuthorAlreadyDeleted = ErrorCode{"05", "Author already deleted"}

	InvalidRequest = ErrorCode{"98", "Invalid request"}
)
