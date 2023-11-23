package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// That handles user login.
// It parses the request body, validates the input using a validator, queries the database for the user,
// verifies the provided password against the stored password hash, and generates a JWT token upon successful login.
// It returns a JSON response with the generated token.
//
// Usage:
//
//	app.Post("/login", s.login)
func (s *Server) login(c *fiber.Ctx) error {
	body := new(LoginBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res := s.users.FindOne(c.Context(), bson.D{{Key: "username", Value: body.Username}})
	user := &entities.User{}

	res.Decode(user)

	err = auth.Verify([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, "invalid password")
	}

	token, err := auth.Sign(&auth.Payload{Id: user.Id}, []byte(os.Getenv("JWT_SECRET")), nil)
	if err != nil {
		log.Println(err)
		return errors.NewHttpError(c, errors.INTERNAL_SERVER_ERROR, "there is an error while signing your token")
	}

	return c.JSON(LoginResponse{Token: token})
}
