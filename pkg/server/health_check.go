package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type healthCheckResponse struct {
	Running    bool      `json:"running"`
	ServerTime time.Time `json:"serverTime"`
}

func (s *Server) healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(healthCheckResponse{
		Running:    true,
		ServerTime: time.Now().Local(),
	})
}
