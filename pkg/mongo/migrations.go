package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/do"
	"github.com/telle-bots/bot-runner/pkg/config"
	"github.com/telle-bots/bot-runner/pkg/mongo/migrations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Migration struct {
	Name       string                                              `bson:"_id"`
	Run        func(ctx context.Context, db *mongo.Database) error `bson:"-"`
	ExecutedAt time.Time                                           `bson:"executed_at"`
}

type Migrator struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	db         *mongo.Database
	coll       *mongo.Collection
	done       <-chan struct{}
	migrations []Migration
}

func NewMigrator(in *do.Injector) (*Migrator, error) {
	db, err := do.Invoke[*mongo.Database](in)
	if err != nil {
		return nil, fmt.Errorf("mongo database: %w", err)
	}

	return &Migrator{
		cfg:  do.MustInvoke[*config.Config](in),
		log:  do.MustInvoke[*zap.SugaredLogger](in),
		db:   db,
		coll: db.Collection("migrations"),
		migrations: []Migration{
			{
				Name: "bot-token-index",
				Run:  migrations.BotTokenIndex,
			},
		},
	}, nil
}

func (m *Migrator) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.MongoMigrationTimeout)
	defer cancel()

	m.done = ctx.Done()

	for _, migration := range m.migrations {
		err := m.coll.FindOne(ctx, bson.D{{"_id", migration.Name}}).Err()
		if err == nil {
			m.log.Infow("Skipping migration", "name", migration.Name)
			continue
		}

		if !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("get migration: %s, err: %w", migration.Name, err)
		}

		m.log.Infow("Running migration", "name", migration.Name)

		if err = migration.Run(ctx, m.db); err != nil {
			return fmt.Errorf("migration failed: %s, err: %w", migration.Name, err)
		}

		migration.ExecutedAt = time.Now()
		if _, err = m.coll.InsertOne(ctx, migration); err != nil {
			return fmt.Errorf("save migration run: %s, err: %w", migration.Name, err)
		}

		m.log.Infow("Migration done", "name", migration.Name)
	}

	return nil
}

func (m *Migrator) Shutdown() error {
	<-m.done
	return nil
}
