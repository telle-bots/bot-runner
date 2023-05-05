package mongo

import (
	"fmt"

	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func Database(in *do.Injector) (*mongo.Database, error) {
	cfg := do.MustInvoke[*config.Config](in)

	client, err := do.Invoke[*Client](in)
	if err != nil {
		return nil, fmt.Errorf("mongo clinent: %w", err)
	}

	return client.Database(cfg.MongoDatabase), nil
}
