package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MacPiggins/gw-notification/internal/app"
	"github.com/MacPiggins/gw-notification/internal/config"
	"github.com/MacPiggins/gw-notification/internal/logging"
)

func main() {
	h := &logging.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
	slog.SetDefault(slog.New(h))
	slog.SetLogLoggerLevel(slog.LevelInfo)

	conf, err := config.Load("config.env")
	if err != nil {
		slog.Error("error while loading config", slog.Any("error", err))
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err = app.Run(ctx, conf)
	if err != nil {
		slog.Error("error while running notification service", slog.Any("error", err))
		return
	}
}
