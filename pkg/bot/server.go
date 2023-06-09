package bot

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const RunnerTask = "bot-runner"
const RunnerQueue = "default"

const delimiter = ":"

type TaskPayload struct {
	BotID int64 `json:"bot_id"`
}

func NewBotTask(botID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(TaskPayload{BotID: botID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(RunnerTask, payload), nil
}

type TaskServer struct {
	cfg              *config.Config
	log              *zap.SugaredLogger
	botRepository    *Repository
	botTaskInspector *TaskInspector
	srv              *asynq.Server
	mux              *asynq.ServeMux
	webhookSrv       *WebhookServer
}

func NewTaskServer(in *do.Injector) (*TaskServer, error) {
	cfg := do.MustInvoke[*config.Config](in)
	log := do.MustInvoke[*zap.SugaredLogger](in)

	botRepository, err := do.Invoke[*Repository](in)
	if err != nil {
		return nil, fmt.Errorf("bot repository: %w", err)
	}

	redisCfg, err := do.Invoke[*asynq.RedisClientOpt](in)
	if err != nil {
		return nil, fmt.Errorf("redis config: %w", err)
	}

	botTaskInspector, err := do.Invoke[*TaskInspector](in)
	if err != nil {
		return nil, fmt.Errorf("bot task inspector: %w", err)
	}

	webhookSrv, err := do.Invoke[*WebhookServer](in)
	if err != nil {
		return nil, fmt.Errorf("webhook server: %w", err)
	}

	srv := asynq.NewServer(
		redisCfg,
		asynq.Config{
			Concurrency:     cfg.MaxTasks,
			ShutdownTimeout: cfg.TaskShutdownTimeout,
			Logger:          log.WithOptions(zap.AddCallerSkip(1)),
		},
	)
	mux := asynq.NewServeMux()

	runner := &TaskServer{
		cfg:              cfg,
		log:              do.MustInvoke[*zap.SugaredLogger](in),
		botRepository:    botRepository,
		botTaskInspector: botTaskInspector,
		srv:              srv,
		mux:              mux,
		webhookSrv:       webhookSrv,
	}
	runner.init()

	return runner, nil
}

func (s *TaskServer) init() {
	s.mux.HandleFunc(RunnerTask, s.runBotTask)
}

func (s *TaskServer) Start() error {
	return s.srv.Start(s.mux)
}

func (s *TaskServer) Shutdown() error {
	s.srv.Stop()
	s.srv.Shutdown()
	return nil
}

func (s *TaskServer) runBotTask(ctx context.Context, t *asynq.Task) error {
	if t.Type() != RunnerTask {
		return asynq.SkipRetry
	}

	var payload TaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return errors.Join(err, asynq.SkipRetry)
	}

	botInfo, err := s.botRepository.BotInfo(ctx, payload.BotID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Join(err, asynq.SkipRetry)
		}
		return fmt.Errorf("get bot info, id: %d: %w", payload.BotID, err)
	}

	if !botInfo.Enabled {
		return nil
	}
	taskID, _ := asynq.GetTaskID(ctx)

	bot, err := telego.NewBot(botInfo.Token)
	if err != nil {
		return errors.Join(err, asynq.SkipRetry)
	}

	me, err := bot.GetMe()
	if err != nil {
		return err
	}
	s.log.Infow("Starting bot", "bot-id", me.ID, "username", me.Username)
	defer s.log.Infow("Stopped bot", "bot-id", me.ID, "username", me.Username)

	tokenData := sha256.Sum256([]byte(taskID + delimiter + botInfo.Token))
	secretToken := hex.EncodeToString(tokenData[:])

	updates, err := bot.UpdatesViaWebhook(
		strconv.FormatInt(botInfo.ID, 10)+delimiter+secretToken,
		telego.WithWebhookContext(ctx),
		telego.WithWebhookSet(&telego.SetWebhookParams{
			URL:            fmt.Sprintf("%s/bots/%d/webhook", s.cfg.WebhookBaseURL, botInfo.ID),
			AllowedUpdates: []string{telego.MessageUpdates},
			SecretToken:    secretToken,
		}),
		telego.WithWebhookServer(&telego.NoOpBotWebhookServer{
			RegisterHandlerFunc: s.webhookSrv.RegisterBot,
		}),
	)
	if err != nil {
		return err
	}

	bh, err := th.NewBotHandler(bot, updates, th.WithStopTimeout(s.cfg.BotHandlerShutdownTimeout))
	if err != nil {
		return err
	}

	// var wfs logic_2.Workflows

	// ac := logic.Actions(actions.NewActionBot(bot))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// logic_2.ProcessWorkflows(
		// 	s.log,
		// 	wfs,
		// 	[]logic_2.TriggerType{logic_2.TriggerTypeBot},
		// 	[]logic_2.EventType{logic_2.EventTypeUpdate},
		// 	update,
		// )

		/*
			tr := logic.Triggers(triggers.NewTriggerUpdate(update))

			trigger, ok := tr[logic.TriggerMessageText]
			if !ok {
				return
			}

			if ok, err = trigger(triggers.TriggerArgs{
				Condition: triggers.MessageTextCondition{
					Text: conditions.ConditionString{
						StartsWith: telego.ToPtr("Send: "),
					},
				},
			}); err != nil {
				s.log.Error(err)
				return
			} else if !ok {
				return
			}

			action, ok := ac[logic.ActionSendMessage]
			if !ok {
				return
			}

			_, err = action(actions.ActionArgs{
				Data: actions.SendMessageData{
					Text: strings.TrimPrefix(update.Message.Text, "Send: "),
				},
				Context: actions.SendMessageContext{
					ChatID: update.Message.Chat.ID,
				},
			})
			if err != nil {
				s.log.Error(err)
				return
			}
		*/
	})

	go bh.Start()

	<-ctx.Done()

	_ = bot.StopWebhook()
	bh.Stop()

	s.webhookSrv.UnregisterBot(botInfo.ID)

	newCtx := context.Background()
	botInfo, err = s.botRepository.BotInfo(newCtx, payload.BotID)
	if err != nil {
		return fmt.Errorf("get bot info, id: %d: %w", payload.BotID, err)
	}

	if botInfo.Enabled {
		return nil
	}

	return s.botTaskInspector.DeleteTask(RunnerQueue, taskID)
}
