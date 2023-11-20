package server

import (
	"go-api/pkg/entities"
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePostBody struct {
	Title       string `json:"title" validate:"required,alpha,min=2,max=30"`
	Description string `json:"description" validate:"required,min=20,max=200"`
}

type CreatePostResponse struct {
	Post *entities.Post `json:"post"`
}

func (s *Server) createPost(c *fiber.Ctx) error {
	body := new(CreatePostBody)
	if err := c.BodyParser(body); err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	err := s.validator.Struct(body)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
	}

	post := new(entities.Post)

	post.Title = body.Title
	post.Description = body.Description
	post.Description = body.Description
	post.Status = entities.PENDING

	res, err := s.posts.InsertOne(c.Context(), post)
	if err != nil {
		return errors.NewHttpError(c, errors.BAD_REQUEST, "failed to create post")
	}

	post.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return c.JSON(CreatePostResponse{post})
}
