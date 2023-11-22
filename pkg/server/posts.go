package server

import (
	"go-api/pkg/entities"
	"go-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostsResponse struct {
	Posts []entities.Post `json:"posts"`
}

// That retrieves posts with the specified status.
// It queries the database for posts with the given status and sends a JSON response with the retrieved posts.
//
// Usage:
//
//	app.Get("/posts/:status", s.getPosts(entities.PostStatus))
func (s *Server) getPosts(status entities.PostStatus) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemsPerPage := c.QueryInt("itemsPerPage", 10)
		page := c.QueryInt("page", 1)
		skip := int64(page*itemsPerPage - itemsPerPage)
		limit := int64(itemsPerPage)

		res, err := s.posts.Find(c.Context(), bson.D{{Key: "status", Value: status}}, &options.FindOptions{Skip: &skip, Limit: &limit})
		if err != nil {
			return errors.NewHttpError(c, errors.BAD_REQUEST, err.Error())
		}

		posts := []entities.Post{}
		if err := res.All(c.Context(), &posts); err != nil {
			return errors.NewHttpError(c, errors.INTERNAL_SERVER_ERROR, err.Error())
		}

		return c.JSON(PostsResponse{posts})
	}
}
