package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/telle-bots/bot-runner/pkg/logic_1"
)

func (s *Server) actions(ctx *fiber.Ctx) error {
	return ctx.JSON(logic_1.ActionDefinitions)
}

func (s *Server) triggers(ctx *fiber.Ctx) error {
	return ctx.JSON(logic_1.TriggerDefinitions)
}
