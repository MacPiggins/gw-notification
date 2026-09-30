package storage

import (
	"context"
	"testing"
	"time"

	"github.com/MacPiggins/gw-notification/internal/model"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestNewMongo(t *testing.T) {
	db, err := NewMongo(context.Background(), startMongo(t))
	if err != nil {
		t.Fatalf("NewMongo() error = %v", err)
	}
	if db == nil || db.client == nil {
		t.Fatal("NewMongo() returned an empty storage")
	}

	if err := db.Disconnect(context.Background()); err != nil {
		t.Fatalf("Disconnect() error = %v", err)
	}
}

func TestMongoSaveTransfer(t *testing.T) {
	uri := startMongo(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := NewMongo(ctx, uri)
	if err != nil {
		t.Fatalf("NewMongo() error = %v", err)
	}
	defer db.Disconnect(context.Background())

	transfer := &model.Transfer{
		TransactionID: "tx-123",
		Sender:        "alice",
		Receiver:      "bob",
		Amount:        100,
		Currency:      "USD",
		OccurredAt:    time.Now().UTC(),
	}
	if err := db.SaveTransfer(ctx, transfer); err != nil {
		t.Fatalf("SaveTransfer() error = %v", err)
	}
}

func startMongo(t *testing.T) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	container, err := mongodb.Run(ctx, "mongo:7", mongodb.WithUsername("test"), mongodb.WithPassword("test"))
	if err != nil {
		t.Fatalf("starting MongoDB container: %v", err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Errorf("terminating MongoDB container: %v", err)
		}
	})

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("getting MongoDB connection string: %v", err)
	}

	return uri
}
