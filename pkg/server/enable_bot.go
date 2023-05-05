package server

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/telle-bots/bot-runner/pkg/bot"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) enableBot(ctx *fiber.Ctx) error {
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

	if botInfo.Enabled {
		return ctx.SendStatus(fiber.StatusOK)
	}

	task, err := bot.NewBotTask(botID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	taskInfo, err := s.botTaskClient.EnqueueContext(c, task,
		asynq.Queue(bot.RunnerQueue), asynq.Timeout(0), asynq.Unique(s.cfg.TaskTTL))
	// TODO: Check case when botInfo.TaskID == ""
	if errors.Is(err, asynq.ErrDuplicateTask) && botInfo.TaskID != "" {
		taskInfo = &asynq.TaskInfo{ID: botInfo.TaskID}
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	botInfo.Enabled = true
	botInfo.TaskID = taskInfo.ID
	if err = s.botRepository.UpdateBotInfo(c, botInfo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error()) // FIXME
	}

	return ctx.SendStatus(fiber.StatusOK)
}
