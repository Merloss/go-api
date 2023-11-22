package server

import (
	"github.com/gofiber/fiber/v2"
)

// That indicates the liveness of the server.
//
// Usage:
//
//	app.Get("/liveness", s.liveness)
func (s *Server) liveness(c *fiber.Ctx) error {
	return c.SendString("OK")
}
