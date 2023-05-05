package main

import (
	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/cmd"
	"github.com/telle-bots/bot-runner/pkg/mongo"
	"github.com/telle-bots/bot-runner/pkg/server"
	"go.uber.org/zap"
)

func main() {
	cmd.Run(func(in *do.Injector, log *zap.SugaredLogger) {
		migrator, err := do.Invoke[*mongo.Migrator](in)
		if err != nil {
			log.Fatalw("Create migrator", "error", err)
		}

		webSrv, err := do.Invoke[*server.Server](in)
		if err != nil {
			log.Fatalw("Create web server", "error", err)
		}

		if err = migrator.Run(); err != nil {
			log.Fatalw("Run migrations", "error", err)
		}

		go func() {
			if err = webSrv.Start(); err != nil {
				log.Fatalw("Start web server", "error", err)
			}
		}()
	})
}
