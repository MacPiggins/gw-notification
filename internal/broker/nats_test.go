package broker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/MacPiggins/gw-notification/internal/config"
	"github.com/MacPiggins/gw-notification/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/testcontainers/testcontainers-go"
	natscontainer "github.com/testcontainers/testcontainers-go/modules/nats"
)

func TestNatsGetMessage(t *testing.T) {
	t.Log("starting NATS integration test")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	t.Log("starting NATS container")
	container, err := natscontainer.Run(ctx, "nats:2.10-alpine", testcontainers.WithCmd("-js"))
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background()) //nolint:errcheck

	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("NATS container started: %s", connectionString)
	conn, err := nats.Connect(connectionString)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	t.Log("connected to NATS")

	js, err := jetstream.New(conn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "NOTIFICATIONS",
		Subjects: []string{"wallet.transfers"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err = js.CreateOrUpdateConsumer(ctx, "NOTIFICATIONS", jetstream.ConsumerConfig{
		Name:          "NOTIFICATION_CONSUMER",
		FilterSubject: "wallet.transfers",
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("created stream and consumer")

	br, err := NewNats(ctx, &config.Config{
		NatsConnstr:  connectionString,
		NatsStream:   "NOTIFICATIONS",
		NatsConsumer: "NOTIFICATION_CONSUMER",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer br.conn.Close()
	t.Log("created NATS broker")

	transfer := model.Transfer{
		TransactionID: "tx-123",
		Sender:        "alice",
		Receiver:      "bob",
		Amount:        100,
		Currency:      "USD",
		OccurredAt:    time.Now().UTC(),
	}
	payload, err := json.Marshal(transfer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = js.Publish(ctx, "wallet.transfers", payload); err != nil {
		t.Fatal(err)
	}
	t.Log("published transfer")

	received, _, err := br.GetMessage()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("received notification")
	if received == nil {
		t.Fatal("expected a transfer")
	}
	if received.ReceivedAt.IsZero() {
		t.Fatal("expected ReceivedAt to be set")
	}
	if received.TransactionID != transfer.TransactionID {
		t.Fatalf("expected transaction ID %q, got %q", transfer.TransactionID, received.TransactionID)
	}
}
