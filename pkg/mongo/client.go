package mongo

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
}

func ConnectClient(in *do.Injector) (*Client, error) {
	cfg := do.MustInvoke[*config.Config](in)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoConnectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConnection))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) Shutdown() error {
	return c.Disconnect(nil)
}
