package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/di"
	"go.uber.org/zap"
)

func Run(run func(in *do.Injector, log *zap.SugaredLogger)) {
	in, err := di.Init()
	assert(err == nil, "init: ", err)

	log := do.MustInvoke[*zap.SugaredLogger](in)
	log.Info("Starting...")

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Info("Stopping...")

		err = in.Shutdown()
		if err != nil {
			log.Fatalw("Shutdown", "error", err)
		}

		done <- struct{}{}
	}()

	run(in, log)

	<-done
	log.Info("Done")
	_ = log.Sync()
}
