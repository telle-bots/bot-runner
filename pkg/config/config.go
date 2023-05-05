package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/samber/do"
)

const envPrefix = "BOT_RUNNER_"

type Config struct {
	ManagerListenAddress string        `env:"MANAGER_LISTEN_ADDRESS,required" validate:"hostname_port"`
	RunnerListenAddress  string        `env:"RUNNER_LISTEN_ADDRESS,required" validate:"hostname_port"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT,required" validate:"gte=0"`

	LoggerType string `env:"LOGGER_TYPE,required" validate:"oneof=dev prod"`

	MongoConnection        string        `env:"MONGO_CONNECTION,required"`
	MongoDatabase          string        `env:"MONGO_DATABASE,required"`
	MongoConnectionTimeout time.Duration `env:"MONGO_CONNECTION_TIMEOUT,required" validate:"gte=0"`
	MongoMigrationTimeout  time.Duration `env:"MONGO_MIGRATION_TIMEOUT,required" validate:"gte=0"`

	RedisAddress              string        `env:"REDIS_ADDRESS,required" validate:"hostname_port"`
	TaskTTL                   time.Duration `env:"TASK_TTL,required" validate:"gte=1s"`
	MaxTasks                  int           `env:"MAX_TASKS,required" validate:"gt=0"`
	TaskShutdownTimeout       time.Duration `env:"TASK_SHUTDOWN_TIMEOUT,required" validate:"gt=0"`
	WebhookBaseURL            string        `env:"WEBHOOK_BASE_URL,required" validate:"url"`
	BotHandlerShutdownTimeout time.Duration `env:"BOT_HANDLER_SHUTDOWN_TIMEOUT,required" validate:"gte=0"`
}

func LoadConfig(in *do.Injector) (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg, env.Options{
		Prefix: envPrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	validate, err := do.Invoke[*validator.Validate](in)
	if err != nil {
		return nil, fmt.Errorf("validator: %w", err)
	}

	if err = validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &cfg, nil
}
