package bot

import (
	"context"
	"fmt"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(in *do.Injector) (*Repository, error) {
	db, err := do.Invoke[*mongo.Database](in)
	if err != nil {
		return nil, fmt.Errorf("mongo databse: %w", err)
	}

	return &Repository{
		coll: db.Collection("bots"),
	}, nil
}

func (r *Repository) CreateBot(ctx context.Context, id int64, token string) error {
	_, err := r.coll.InsertOne(ctx, Info{
		ID:    id,
		Token: token,
	})
	return err
}

func (r *Repository) Bots(ctx context.Context) ([]Info, error) {
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var bots []Info
	if err = cursor.All(ctx, &bots); err != nil {
		return nil, err
	}

	return bots, nil
}

func (r *Repository) BotInfo(ctx context.Context, id int64) (Info, error) {
	resp := r.coll.FindOne(ctx, bson.D{{"_id", id}})

	var bot Info
	if err := resp.Decode(&bot); err != nil {
		return Info{}, err
	}

	return bot, nil
}

func (r *Repository) UpdateBotInfo(ctx context.Context, info Info) error {
	_, err := r.coll.UpdateByID(ctx, info.ID, bson.D{{"$set", info}})
	return err
}
