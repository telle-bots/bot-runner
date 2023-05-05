package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/telego"
)

type newBotRequest struct {
	Token string `json:"token" validate:"required"`
}

func (s *Server) newBot(ctx *fiber.Ctx) error {
	var request newBotRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	if err := s.validate.Struct(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error()) // FIXME
	}

	bot, err := telego.NewBot(request.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error()) // FIXME
	}

	botUser, err := bot.GetMe()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error()) // FIXME
	}

	if err = s.botRepository.CreateBot(ctx.Context(), botUser.ID, request.Token); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	return ctx.JSON(botUser.ID)
}
