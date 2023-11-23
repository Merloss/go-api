package server

import (
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// That deletes a post based on the provided post ID.
// It parses the post ID from the request parameters, converts it to a primitive.ObjectID, and deletes the post from the database.
// It returns a success message if the deletion is successful.
//
// Usage:
//
//	app.Delete("/posts/:id", s.deletePost)
func (s *Server) deletePost(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}
	err = s.posts.FindOneAndDelete(c.Context(), bson.D{{Key: "_id", Value: oid}}).Err()

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	return c.SendString("success")
}
