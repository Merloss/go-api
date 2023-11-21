package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/errors"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) roleGuard(roles ...auth.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userId := claims["id"].(string)

		oid, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
		}

		res := s.users.FindOne(c.Context(), bson.D{{Key: "_id", Value: oid}})
		user := &entities.User{}
		res.Decode(user)

		if slices.Contains(user.Roles, auth.ADMIN) {
			return c.Next()
		}

		if !HasAllFields(user.Roles, roles) {
			return errors.NewHttpError(c, errors.BAD_REQUEST, "Insufficient permission")
		}

		c.Locals("user", user)
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
