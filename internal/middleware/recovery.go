package middleware

import (
	"book-manager/internal/dto/common"
	"book-manager/internal/utils/logger"
	"fmt"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

// Recovery middleware recovers from panics and returns a proper error response
func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					requestID := GetRequestID(c)

					// Log the panic
					err := fmt.Errorf("%v", r)
					stack := string(debug.Stack())
					logger.Error("PANIC", nil, "[%s] Panic recovered: %v\n%s", requestID, err, stack)

					// Return error response
					resp := common.Response[any]{
						ResponseId:   requestID,
						ResponseCode: "99",
						ResponseMsg:  "Internal server error",
						Data:         nil,
					}
					c.JSON(500, resp)
				}
			}()

			return next(c)
		}
	}
}
