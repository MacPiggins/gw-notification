package storage
import (
	"context"
	"log/slog"

	"github.com/MacPiggins/gw-notification/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
}

func NewMongo(ctx context.Context, connstr string) (*Mongo, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connstr))
	if err != nil {
		slog.ErrorContext(ctx, "unable to connect to mongo error", slog.Any("error", err))
		return nil, err
	}
	return &Mongo{client: client}, nil
}

func (db *Mongo) Disconnect(ctx context.Context) error {
	err := db.client.Disconnect(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "unable to disconnect from mongo", slog.Any("error", err))
	}
	return nil
}

func (db *Mongo) SaveTransfer(ctx context.Context, transfer *model.Transfer) error {
	collection := db.client.Database("gw").Collection("transfers")
	_, err := collection.InsertOne(ctx, transfer)
	if err != nil {
		slog.ErrorContext(ctx, "failed to save transfer", slog.Any("error", err))
	}
	return nil
}