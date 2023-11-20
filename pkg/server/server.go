package server

import (
	"context"
	_auth "go-api/pkg/auth"
	"net"
	"os"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	app       *fiber.App
	validator *validator.Validate
	users     *mongo.Collection
	posts     *mongo.Collection
}

func New(db *mongo.Database) *Server {
	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		ContextKey: "requestId",
	}))
	app.Use(logger.New(logger.Config{
		Format:        "${pid} ${locals:requestId} ${status} - ${method} ${path}\n",
		DisableColors: true,
	}))

	validator := validator.New()

	users := db.Collection("users")
	contents := db.Collection("posts")

	s := &Server{app, validator, users, contents}
	s.init()

	return s
}

func (s *Server) init() {
	api := s.app.Group("/api")
	auth := api.Group("/auth")
	posts := api.Group("/posts")

	jwtware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	})

	api.Get("/liveness", s.liveness)

	auth.Post("/register", s.register)
	auth.Post("/login", s.login)

	posts.Get("/", jwtware, s.roleGuard(_auth.VIEWER), s.getPosts)
	posts.Post("/", jwtware, s.roleGuard(_auth.EDITOR), s.createPost)
	posts.Patch("/:id", jwtware, s.roleGuard(_auth.EDITOR), s.editPost)
	posts.Delete("/:id", jwtware, s.roleGuard(_auth.ADMIN), s.deletePost)
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(net.JoinHostPort("127.0.0.1", port))
}

func (s *Server) Close(ctx context.Context) error {
	return s.app.Server().ShutdownWithContext(ctx)
}
