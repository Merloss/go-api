package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/errors"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) roleGuard(role ...auth.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		r := claims["role"].(auth.Role)

		if r == auth.ADMIN {
			return c.Next()
		}

		if !slices.Contains(role, r) {
			return errors.NewHttpError(c, errors.BAD_REQUEST, "Insufficient permission")
		}
		return c.Next()
	}
}
