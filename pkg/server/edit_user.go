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
	Id       string      `json:"id,omitempty"`
	Roles    []auth.Role `json:"roles"`
	Username string      `json:"username"`
}

type UpdateUserResponse struct {
	User UpdateUserBody `json:"user"`
}

func (s *Server) editUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	body := new(UpdateUserBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	user := &entities.User{}

	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res := s.users.FindOneAndUpdate(c.Context(), bson.D{{Key: "_id", Value: oid}}, bson.D{{Key: "$set", Value: bson.D{{Key: "roles", Value: body.Roles}, {Key: "username", Value: body.Username}}}})

	err = res.Err()
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res.Decode(user)

	return c.JSON(UpdateUserResponse{*body})
}
