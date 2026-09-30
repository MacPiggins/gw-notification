package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	StorageConnstr string
	NatsConnstr    string
	NatsStream     string
	NatsConsumer   string
}

const (
	defaultStorageConnstr = "http://localhost:9000"
	defaultNatsConnstr    = "nats://localhost:4222"
	defaultNatsStream     = "wallet"
	defaultNatsConsumer   = "notification"
)

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func Load(files ...string) (*Config, error) {
	if err := godotenv.Load(files...); err != nil {
		return nil, err
	}
	return &Config{
		StorageConnstr: envOrDefault("STORAGECONNSTR", defaultStorageConnstr),
		NatsConnstr:    envOrDefault("NATSCONNSTR", defaultNatsConnstr),
		NatsStream:     envOrDefault("NATSSTREAM", defaultNatsStream),
		NatsConsumer:   envOrDefault("NATSCONSUMER", defaultNatsConsumer),
		},
		nil
}
