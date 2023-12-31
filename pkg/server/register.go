package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

// That handles user registration.
// It parses the request body, validates the input using a validator, checks if the username is already taken,
// creates a new user in the database with hashed password, and generates a JWT token upon successful registration.
// It returns a JSON response with the created user and the generated token.
//
// Usage:
//
//	app.Post("/register", s.register)
func (s *Server) register(c *fiber.Ctx) error {
	body := new(RegisterBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	existUserCount, err := s.users.CountDocuments(c.Context(), bson.D{{Key: "username", Value: body.Username}})

	if existUserCount > 0 {
		return errors.NewHttpError(c, errors.BAD_REQUEST, "username is already taken")
	}

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	user := new(entities.User)
	user.Username = body.Username
	user.Password = string(auth.Hash(body.Password))
	user.Roles = []auth.Role{auth.VIEWER}

	res, err := s.users.InsertOne(c.Context(), user)
	if err != nil {
		return errors.NewHttpError(c, errors.INTERNAL_SERVER_ERROR, "error")
	}

	user.Id = res.InsertedID.(primitive.ObjectID).Hex()

	token, err := auth.Sign(&auth.Payload{Id: user.Id}, []byte(os.Getenv("JWT_SECRET")), nil)
	if err != nil {
		log.Println(err)
		return errors.NewHttpError(c, errors.INTERNAL_SERVER_ERROR, "there is an error while signing your token")
	}

	return c.JSON(RegisterResponse{User: user, Token: token})
}
