package bot

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/telego"
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.uber.org/zap"
)

type botContext struct {
	handler     telego.WebhookHandler
	secretToken string
}

type WebhookServer struct {
	cfg *config.Config
	log *zap.SugaredLogger
	app *fiber.App

	lock        *sync.Mutex
	botHandlers map[int64]botContext
}

func NewWebhookServer(in *do.Injector) (*WebhookServer, error) {
	return &WebhookServer{
		cfg: do.MustInvoke[*config.Config](in),
		log: do.MustInvoke[*zap.SugaredLogger](in),
		app: fiber.New(fiber.Config{
			AppName:               "Bot Webhook Server",
			DisableStartupMessage: true,
		}),

		lock:        &sync.Mutex{},
		botHandlers: map[int64]botContext{},
	}, nil
}

func (w *WebhookServer) Start() error {
	w.app.Post("/bots/:botID<int>/webhook", w.botHandler)
	return w.app.Listen(w.cfg.RunnerListenAddress)
}

func (w *WebhookServer) Shutdown() error {
	return w.app.ShutdownWithTimeout(w.cfg.ShutdownTimeout)
}

func (w *WebhookServer) RegisterBot(path string, handler telego.WebhookHandler) error {
	botIDStr, token, ok := strings.Cut(path, delimiter)
	if !ok {
		return fmt.Errorf("invalid bot handler path: %q", path)
	}

	botID, _ := strconv.ParseInt(botIDStr, 10, 64)

	w.lock.Lock()
	defer w.lock.Unlock()

	w.botHandlers[botID] = botContext{
		handler:     handler,
		secretToken: token,
	}

	return nil
}

func (w *WebhookServer) UnregisterBot(botID int64) {
	w.lock.Lock()
	defer w.lock.Unlock()

	delete(w.botHandlers, botID)
}

func (w *WebhookServer) botHandler(ctx *fiber.Ctx) error {
	botID, _ := strconv.ParseInt(ctx.Params("botID"), 10, 64)

	bot, ok := w.botHandlers[botID]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	token := ctx.Get(telego.WebhookSecretTokenHeader)
	if token != bot.secretToken {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	if err := bot.handler(ctx.Body()); err != nil {
		w.log.Errorw("webhook bot handler", "bot-id", botID, "error", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
