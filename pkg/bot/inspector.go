package bot

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/samber/do"
)

type TaskInspector struct {
	*asynq.Inspector
}

func NewTaskInspector(in *do.Injector) (*TaskInspector, error) {
	redisCfg, err := do.Invoke[*asynq.RedisClientOpt](in)
	if err != nil {
		return nil, fmt.Errorf("redis config: %w", err)
	}

	return &TaskInspector{
		Inspector: asynq.NewInspector(redisCfg),
	}, nil
}

func (i *TaskInspector) Shutdown() error {
	return i.Close()
}
