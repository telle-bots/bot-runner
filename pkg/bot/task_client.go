package bot

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/samber/do"
)

type TaskClient struct {
	*asynq.Client
}

func NewTaskClient(in *do.Injector) (*TaskClient, error) {
	redisCfg, err := do.Invoke[*asynq.RedisClientOpt](in)
	if err != nil {
		return nil, fmt.Errorf("redis config: %w", err)
	}

	return &TaskClient{
		Client: asynq.NewClient(redisCfg),
	}, nil
}

func (c *TaskClient) Shutdown() error {
	return c.Close()
}
