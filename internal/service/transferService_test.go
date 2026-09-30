package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MacPiggins/gw-notification/internal/cache/inmemory"
	"github.com/MacPiggins/gw-notification/internal/model"
)

type transferStorageStub struct {
	saved *model.Transfer
	err   error
}

func (s *transferStorageStub) SaveTransfer(_ context.Context, transfer *model.Transfer) error {
	s.saved = transfer
	return s.err
}

type brokerStub struct {
	transfer *model.Transfer
	err      error
}

func (b *brokerStub) GetMessage() (*model.Transfer, func(), error) {
	return b.transfer, nil, b.err
}

func TestTransferServiceProcessMessage(t *testing.T) {
	transfer := &model.Transfer{
		TransactionID: "tx-123",
		Sender:        "alice",
		Receiver:      "bob",
		Amount:        100,
		Currency:      "USD",
		OccurredAt:    time.Now().UTC(),
	}
	brokerErr := errors.New("get message failed")
	storageErr := errors.New("save transfer failed")

	tests := []struct {
		name       string
		brokerErr  error
		storageErr error
		wantErr    error
		wantSaved  bool
	}{
		{name: "success", wantSaved: true},
		{name: "broker error", brokerErr: brokerErr, wantErr: brokerErr},
		{name: "storage error", storageErr: storageErr, wantErr: storageErr, wantSaved: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &transferStorageStub{err: tt.storageErr}
			broker := &brokerStub{transfer: transfer, err: tt.brokerErr}
			cache := inmemory.NewCache()
			svc := NewTransferService(storage, broker, cache)

			if err := svc.processTransfer(context.Background()); !errors.Is(err, tt.wantErr) {
				t.Fatalf("processTransfer() error = %v, want %v", err, tt.wantErr)
			}
			if (storage.saved != nil) != tt.wantSaved {
				t.Fatalf("saved = %v, want %v", storage.saved != nil, tt.wantSaved)
			}
			if tt.wantSaved && storage.saved != transfer {
				t.Fatalf("saved transfer = %p, want %p", storage.saved, transfer)
			}
		})
	}
}

func TestTransferServiceStartStopsWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cache := inmemory.NewCache()
	svc := NewTransferService(&transferStorageStub{}, &brokerStub{}, cache)
	if err := svc.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
}
