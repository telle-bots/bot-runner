package logger

import (
	"fmt"

	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.uber.org/zap"
)

func NewLogger(in *do.Injector) (*zap.SugaredLogger, error) {
	cfg := do.MustInvoke[*config.Config](in)

	var err error
	var logger *zap.Logger
	switch cfg.LoggerType {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		err = fmt.Errorf("unknown type: %q", cfg.LoggerType)
	}

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
