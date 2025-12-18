package handler

import (
	"searchav/internal/config"

	"github.com/gofiber/fiber/v2"
)

const (
	// AuthHeader is the header name for password authentication
	AuthHeader = "X-Auth-Password"
)

// AuthMiddleware creates a password authentication middleware
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip if auth is disabled
		if !cfg.Auth.Enabled {
			return c.Next()
		}

		password := c.Get(AuthHeader)
		if !cfg.ValidatePassword(password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code": 401,
				"msg":  "unauthorized",
			})
		}

		return c.Next()
	}
}
