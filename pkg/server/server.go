package server

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/bot"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	cfg              *config.Config
	log              *zap.SugaredLogger
	app              *fiber.App
	validate         *validator.Validate
	botRepository    *bot.Repository
	botTaskClient    *bot.TaskClient
	botTaskInspector *bot.TaskInspector
}

func NewServer(in *do.Injector) (*Server, error) {
	cfg := do.MustInvoke[*config.Config](in)
	log := do.MustInvoke[*zap.SugaredLogger](in)
	validate := do.MustInvoke[*validator.Validate](in)

	botRepository, err := do.Invoke[*bot.Repository](in)
	if err != nil {
		return nil, fmt.Errorf("bot manager: %w", err)
	}

	botTaskClient, err := do.Invoke[*bot.TaskClient](in)
	if err != nil {
		return nil, fmt.Errorf("bot task clinet: %w", err)
	}

	botTaskInspector, err := do.Invoke[*bot.TaskInspector](in)
	if err != nil {
		return nil, fmt.Errorf("bot task inspector: %w", err)
	}

	srv := &Server{
		cfg: cfg,
		log: log,
		app: fiber.New(fiber.Config{
			AppName:               "Bot Manager",
			DisableStartupMessage: true,
		}),
		validate:         validate,
		botRepository:    botRepository,
		botTaskClient:    botTaskClient,
		botTaskInspector: botTaskInspector,
	}
	srv.init()

	return srv, nil
}

func (s *Server) Start() error {
	return s.app.Listen(s.cfg.ManagerListenAddress)
}

func (s *Server) Shutdown() error {
	return s.app.ShutdownWithTimeout(s.cfg.ShutdownTimeout)
}

func (s *Server) init() {
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, panicErr any) {
			s.log.Errorw("Panic", "recover", panicErr, "path", ctx.Path(), "method", ctx.Method())
		},
	}))

	s.app.Get("/health", s.healthCheck)

	s.app.Post("/bots", s.newBot)
	s.app.Get("/bots", s.bots)
	s.app.Post("/bots/:botID/enable", s.enableBot)
	s.app.Post("/bots/:botID/disable", s.disableBot)
}
