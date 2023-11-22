package errors

import (
	"github.com/gofiber/fiber/v2"
)

// Must is a utility function that panics if the provided error is not nil.
// It is commonly used in situations where an error is unexpected and indicates a critical issue.
//
// Usage:
//
//	Must(err)
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// NewHttpError creates a new HTTP error response and sends it to the client using Fiber.
// It constructs an ErrorResponse with the provided error code and message, sets the appropriate HTTP status code,
// and includes the request ID in the response if available in the Fiber context.
//
// Usage:
//
//	err := NewHttpError(c, BAD_REQUEST, "Invalid input")
func NewHttpError(c *fiber.Ctx, code ErrorCode, message string) error {
	res := ErrorResponse{}
	res.Error.Code = code
	res.Error.Message = message

	var status int
	switch code {
	case INTERNAL_SERVER_ERROR:
		status = fiber.StatusInternalServerError
	case NOT_FOUND:
		status = fiber.StatusNotFound
	case BAD_REQUEST:
		status = fiber.StatusBadRequest
	}

	if requestId, ok := c.Locals("requestId").(string); ok {
		res.Error.RequestId = requestId
	}
	return c.Status(status).JSON(res)
}
