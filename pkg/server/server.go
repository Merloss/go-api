package server

import (
	"context"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	app       *fiber.App
	validator *validator.Validate
	users     *mongo.Collection
	contents  *mongo.Collection
}

func New(db *mongo.Database) *Server {
	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		ContextKey: "requestId",
	}))
	app.Use(logger.New(logger.Config{
		Format:        "${pid} ${locals:requestId} ${status} - ${method} ${path}​\n",
		DisableColors: true,
	}))

	validator := validator.New()

	users := db.Collection("users")
	contents := db.Collection("contents")

	s := &Server{app, validator, users, contents}
	s.init()

	return s
}

func (s *Server) init() {
	api := s.app.Group("/api")
	api.Get("/liveness", s.liveness)
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(net.JoinHostPort("127.0.0.1", port))
}

func (s *Server) Close(ctx context.Context) error {

	return s.app.Server().ShutdownWithContext(ctx)
}
