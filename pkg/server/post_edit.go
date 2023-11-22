package server

import (
	"go-api/pkg/auth"
	"go-api/pkg/entities"
	"go-api/pkg/errors"
	"slices"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdatePostBody struct {
	Title       string `json:"title" validate:"required,alpha,min=2,max=30"`
	Description string `json:"description" validate:"required,min=20,max=200"`
	Status      string `json:"status"`
}

type UpdatePostResponse struct {
	Post *entities.Post `json:"post"`
}

// That updates a post based on the provided post ID.
// It retrieves the user information from the Fiber context, parses the post ID from the request parameters,
// parses the request body, validates the input using a validator, and updates the post in the database.
// The update includes changes to the title, description and status fields, with additional permissions for admin users.
// It returns a JSON response with the updated post.
//
// Usage:
//
//	app.Put("/posts/:id", s.editPost)
func (s *Server) editPost(c *fiber.Ctx) error {

	user := c.Locals("user").(*entities.User)
	id := c.Params("id")

	body := new(UpdatePostBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	post := new(entities.Post)

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: body.Title},
				{Key: "description", Value: body.Description},
				{Key: "status", Value: entities.PENDING},
			},
		},
	}

	if slices.Contains(user.Roles, auth.ADMIN) {
		update = bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "title", Value: body.Title},
					{Key: "description", Value: body.Description},
					{Key: "status", Value: body.Status},
				},
			},
		}
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res := s.posts.FindOneAndUpdate(c.Context(), bson.D{{Key: "_id", Value: oid}}, update)

	err = res.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewHttpError(c, errors.NOT_FOUND, err.Error())
		}
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	res.Decode(post)

	return c.JSON(UpdatePostResponse{post})
}
