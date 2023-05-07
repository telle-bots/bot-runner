package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/telle-bots/bot-runner/pkg/logic"
)

func (s *Server) actions(ctx *fiber.Ctx) error {
	return ctx.JSON(logic.ActionDefinitions)
}

func (s *Server) triggers(ctx *fiber.Ctx) error {
	return ctx.JSON(logic.TriggerDefinitions)
}
