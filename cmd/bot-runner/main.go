package main

import (
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/cmd"
	"github.com/telle-bots/bot-runner/pkg/bot"
	"github.com/telle-bots/bot-runner/pkg/mongo"
	"go.uber.org/zap"
)

func main() {
	cmd.Run(func(in *do.Injector, log *zap.SugaredLogger) {
		migrator, err := do.Invoke[*mongo.Migrator](in)
		if err != nil {
			log.Fatalw("Create migrator", "error", err)
		}

		taskSrv, err := do.Invoke[*bot.TaskServer](in)
		if err != nil {
			log.Fatalw("Create task server", "error", err)
		}

		webhookSrv, err := do.Invoke[*bot.WebhookServer](in)
		if err != nil {
			log.Fatalw("Create webhook server", "error", err)
		}

		if err = migrator.Run(); err != nil {
			log.Fatalw("Run migrations", "error", err)
		}

		go func() {
			if err = taskSrv.Start(); err != nil {
				log.Fatalw("Start task server", "error", err)
			}
		}()

		go func() {
			if err = webhookSrv.Start(); err != nil {
				log.Fatalw("Start webhook server", "error", err)
			}
		}()
	})
}
