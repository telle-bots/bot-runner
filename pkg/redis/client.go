package redis

import (
	"github.com/hibiken/asynq"
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
)

func NewRedisClientOpt(in *do.Injector) (*asynq.RedisClientOpt, error) {
	cfg := do.MustInvoke[*config.Config](in)
	return &asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	}, nil
}
