package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/telle-bots/bot-runner/pkg/logic/nodes"
)

func (s *Server) nodes(ctx *fiber.Ctx) error {
	return ctx.JSON(nodes.Nodes)
}
