package server

import (
	"go-api/pkg/entities"
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type PostsResponse struct {
	Posts []entities.Post `json:"posts"`
}

func (s *Server) getPosts(c *fiber.Ctx) error {

	res, err := s.posts.Find(c.Context(), bson.D{{Key: "status", Value: entities.APPROVED}})

	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	posts := []entities.Post{}

	if err := res.All(c.Context(), &posts); err != nil {
		return errors.NewHttpError(c, errors.INTERNAL_SERVER_ERROR, err.Error())
	}

	return c.JSON(PostsResponse{posts})
}
