package handler

import (
	"searchav/internal/constants"
	"searchav/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles global errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := constants.InternalError

	// Handle Fiber errors
	if e, ok := err.(*fiber.Error); ok {
		if e.Code == fiber.StatusBadRequest {
			code = constants.InvalidParams
		}
	}

	return c.Status(fiber.StatusOK).JSON(&dto.Response{
		Code:    code.Code,
		Message: code.Message,
	})
}
