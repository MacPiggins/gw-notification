# GW Notification

A Go service that consumes wallet transfer events from a NATS JetStream stream and stores them in MongoDB.

## How it works

1. Connects to NATS and creates or updates a durable consumer for the configured stream.
2. Reads messages matching the `wallet.transfers` subject and decodes them as transfer events.
3. Adds a `received_at` timestamp and inserts each transfer into the `gw.transfers` MongoDB collection.

The event payload uses this JSON shape:

```json
{
	"transaction_id": "txn-123",
	"sender": "wallet-a",
	"receiver": "wallet-b",
	"amount": 2500,
	"currency": "USD",
	"occurred_at": "2026-09-29T12:00:00Z"
}
```

`amount` is an integer in the smallest unit used by the event producer. `received_at` is assigned by the service and is not read from the incoming JSON.

## Requirements

- Go 1.25 or later
- MongoDB
- NATS Server with JetStream enabled
- Docker, if running the integration tests or using Docker Compose

## Configuration

The service loads `config.env` from its working directory at startup. Create it from the example, then set values for your environment:

```sh
cp config.env.example config.env
```

| Variable | Default | Description |
| --- | --- | --- |
| `STORAGECONNSTR` | `http://localhost:9000` | MongoDB connection URI. Set this to a valid MongoDB URI, such as `mongodb://localhost:27017`. |
| `NATSCONNSTR` | `nats://localhost:4222` | NATS server URL. |
| `NATSSTREAM` | `analytics` | Existing JetStream stream to consume. |
| `NATSCONSUMER` | `analytics` | Durable consumer name created or updated by the service. |

The checked-in storage default is not a valid MongoDB URI; set `STORAGECONNSTR` to your MongoDB URI before starting the service. The database and collection are currently fixed to `gw` and `transfers` in the application.

## Run locally

Start MongoDB and a JetStream-enabled NATS server. Create a stream named `analytics` that captures the `wallet.transfers` subject. For example, with the NATS CLI:

```sh
nats --server nats://localhost:4222 stream add analytics --subjects wallet.transfers --defaults
```

Set `STORAGECONNSTR` and any other overrides in `config.env`, then run:

```sh
go run ./cmd/main.go
```

The process handles interrupt and termination signals. Logs are emitted as JSON to standard output.

## Run with Docker Compose

Docker Compose starts the service, MongoDB, and NATS. The service container connects to dependencies by their Compose service names, so configure `config.env` with container-network addresses before starting it:

```dotenv
STORAGECONNSTR=mongodb://root:password@mongo:27017/?authSource=admin
NATSCONNSTR=nats://nats:4222
NATSSTREAM=analytics
NATSCONSUMER=analytics
```

Start the dependencies first:

```sh
docker compose up -d mongo nats
```

Create the stream from the host (the Compose file publishes NATS on port 4222):

```sh
nats --server nats://localhost:4222 stream add analytics --subjects wallet.transfers --defaults
```

Then build and start the service:

```sh
docker compose up --build -d app
```

To view logs or stop the stack:

```sh
docker compose logs -f app
docker compose down
```

MongoDB data is kept in the `mongo_data` named volume. The Compose file mounts the local `config.env` into the service container.

## Tests

The integration tests use Testcontainers and require a working Docker daemon:

```sh
go test ./...
```

## Current limitations

- The stream must exist before the service starts; the Compose setup does not create it automatically.
- The consumed subject is currently fixed as `wallet.transfers`.
- The current message-processing path does not explicitly acknowledge JetStream messages. Configure and test redelivery behavior before relying on it in production.
- MongoDB database and collection names are currently fixed as `gw.transfers`.
