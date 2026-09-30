package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/MacPiggins/gw-notification/internal/cache/inmemory"
	"github.com/MacPiggins/gw-notification/internal/model"
)

type TransferStorage interface {
	SaveTransfer(ctx context.Context, transfer *model.Transfer) error
}

type Broker interface {
	GetMessage() (*model.Transfer, func(), error)
}

type TransferService struct {
	storage TransferStorage
	broker  Broker
	cache   *inmemory.Cache
}

func NewTransferService(storage TransferStorage, broker Broker, cache *inmemory.Cache) *TransferService {
	return &TransferService{
		storage: storage,
		broker:  broker,
		cache:   cache,
	}
}

func (svc *TransferService) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "transfer service started")
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "transfer service stopped")
			return nil
		default:
			if err := svc.processTransfer(ctx); err != nil {
				slog.ErrorContext(ctx, "failed to process transfer message", "error", err)
				continue
			}
		}
	}
}

func (svc *TransferService) processTransfer(ctx context.Context) error {
	slog.InfoContext(ctx, "waiting for transfer message")
	transfer, ack, err := svc.broker.GetMessage()
	if err != nil {
		slog.ErrorContext(ctx, "failed to retrieve transfer message", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "received transfer message", slog.Any("transfer", transfer))

	if _, ok := svc.cache.Get(transfer.TransactionID); !ok {
		svc.cache.Set(transfer.TransactionID, transfer, time.Hour*24)
	} else {
		slog.WarnContext(ctx, "duplicate transfer message", slog.Any("transfer", transfer))
		return nil
	}
	
	if err := svc.storage.SaveTransfer(ctx, transfer); err != nil {
		slog.ErrorContext(ctx, "failed to save transfer message", slog.Any("error", err))
		return err
	}
	if ack != nil {
		ack()
	}
	return nil
}
