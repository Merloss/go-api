package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/errors"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) roleGuard(roles ...auth.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		r := claims["roles"].([]auth.Role)

		if slices.Contains(r, auth.ADMIN) {
			return c.Next()
		}

		if !HasAllFields(r, roles) {
			return errors.NewHttpError(c, errors.BAD_REQUEST, "Insufficient permission")
		}
		return c.Next()
	}
}

func HasAllFields(userFlags []auth.Role, flags []auth.Role) bool {
	for _, field := range flags {
		found := false
		for _, userFlag := range userFlags {
			if userFlag == field {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
