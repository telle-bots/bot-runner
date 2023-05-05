package server

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) disableBot(ctx *fiber.Ctx) error {
	c := ctx.Context()

	botID, err := strconv.ParseInt(ctx.Params("botID"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error()) // FIXME
	}

	botInfo, err := s.botRepository.BotInfo(c, botID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, err.Error()) // FIXME
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	if !botInfo.Enabled {
		return ctx.SendStatus(fiber.StatusOK)
	}

	if botInfo.TaskID != "" {
		if err = s.botTaskInspector.CancelProcessing(botInfo.TaskID); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
		}
	}

	botInfo.Enabled = false
	botInfo.TaskID = ""
	if err = s.botRepository.UpdateBotInfo(c, botInfo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	return ctx.SendStatus(fiber.StatusOK)
}
