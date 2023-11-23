package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateUserBody struct {
	Roles    []auth.Role `json:"roles,omitempty"`
	Username string      `json:"username,omitempty"`
}

type UpdateUserResponse struct {
	User *entities.User `json:"user"`
}

// Updates a user's information based on the provided request body.
// It parses the user ID from the request parameters, parses the request body, validates the input using a validator,
// and updates the user information in the database. Upon successful update, it returns a JSON response with the updated user.
//
// Usage:
//
//	app.Put("/users/:id", s.editUser)
func (s *Server) editUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	var body UpdateUserBody
	if err := c.BodyParser(&body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(&body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	user := &entities.User{}

	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err = s.users.FindOneAndUpdate(c.Context(), bson.D{{Key: "_id", Value: oid}}, bson.D{{Key: "$set", Value: bson.D{{Key: "roles", Value: body.Roles}, {Key: "username", Value: body.Username}}}}).Decode(user)

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	return c.JSON(UpdateUserResponse{user})
}
