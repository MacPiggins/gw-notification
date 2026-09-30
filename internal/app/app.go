package app

import (
	"context"
	"log/slog"

	"github.com/MacPiggins/gw-notification/internal/broker"
	"github.com/MacPiggins/gw-notification/internal/cache/inmemory"
	"github.com/MacPiggins/gw-notification/internal/config"
	"github.com/MacPiggins/gw-notification/internal/service"
	"github.com/MacPiggins/gw-notification/internal/storage"
)

func Run(ctx context.Context, conf *config.Config) error {
	broker, err := broker.NewNats(ctx, conf)
	if err != nil {
		slog.ErrorContext(ctx, "error initializing nats broker", slog.Any("error", err))
		return err
	}

	storage, err := storage.NewMongo(ctx, conf.StorageConnstr)
	if err != nil {
		slog.ErrorContext(ctx, "error initializing mongo storage", slog.Any("error", err))
		return err
	}
	cache := inmemory.NewCache()
	service := service.NewTransferService(storage, broker, cache)
	if err := service.Start(ctx); err != nil {
		slog.ErrorContext(ctx, "error starting transfer service", slog.Any("error", err))
		return err
	}
	return nil
}
