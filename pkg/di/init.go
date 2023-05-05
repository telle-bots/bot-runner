package di

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/bot"
	"github.com/telle-bots/bot-runner/pkg/config"
	"github.com/telle-bots/bot-runner/pkg/logger"
	"github.com/telle-bots/bot-runner/pkg/mongo"
	"github.com/telle-bots/bot-runner/pkg/redis"
	"github.com/telle-bots/bot-runner/pkg/server"
	"github.com/telle-bots/bot-runner/pkg/validation"
	"go.uber.org/zap"
)

func Init() (*do.Injector, error) {
	in := do.New()

	do.Provide(in, validation.NewValidator)
	do.Provide(in, config.LoadConfig)
	do.Provide(in, logger.NewLogger)
	do.Provide(in, mongo.ConnectClient)
	do.Provide(in, mongo.Database)
	do.Provide(in, mongo.NewMigrator)
	do.Provide(in, redis.NewRedisClientOpt)
	do.Provide(in, server.NewServer)
	do.Provide(in, bot.NewRepository)
	do.Provide(in, bot.NewTaskServer)
	do.Provide(in, bot.NewTaskInspector)
	do.Provide(in, bot.NewTaskClient)
	do.Provide(in, bot.NewWebhookServer)

	if _, err := do.Invoke[*config.Config](in); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if _, err := do.Invoke[*zap.SugaredLogger](in); err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	if _, err := do.Invoke[*validator.Validate](in); err != nil {
		return nil, fmt.Errorf("validator: %w", err)
	}

	return in, nil
}
