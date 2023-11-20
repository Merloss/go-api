package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

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
