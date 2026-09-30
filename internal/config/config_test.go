package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("TEST_CONFIG_VALUE", "configured")
	if got := envOrDefault("TEST_CONFIG_VALUE", "fallback"); got != "configured" {
		t.Fatalf("envOrDefault() = %q, want %q", got, "configured")
	}

	t.Setenv("TEST_CONFIG_VALUE", "")
	if got := envOrDefault("TEST_CONFIG_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("envOrDefault() = %q, want %q", got, "fallback")
	}
}

func TestLoadDefaults(t *testing.T) {
	config, err := Load(writeEnvFile(t, ""))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{
		StorageConnstr: defaultStorageConnstr,
		NatsConnstr:    defaultNatsConnstr,
		NatsStream:     defaultNatsStream,
		NatsConsumer:   defaultNatsConsumer,
	}
	if *config != *want {
		t.Fatalf("Load() = %+v, want %+v", *config, *want)
	}
}

func TestLoadEnvironmentOverrides(t *testing.T) {
	t.Setenv("STORAGECONNSTR", "storage")
	t.Setenv("NATSCONNSTR", "nats")
	t.Setenv("NATSSTREAM", "stream")
	t.Setenv("NATSCONSUMER", "consumer")

	config, err := Load(writeEnvFile(t, ""))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{"storage", "nats", "stream", "consumer"}
	if *config != *want {
		t.Fatalf("Load() = %+v, want %+v", *config, *want)
	}
}

func TestLoadReadsFile(t *testing.T) {
	config, err := Load(writeEnvFile(t, "STORAGECONNSTR=file-storage\nNATSSTREAM=file-stream\n"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.StorageConnstr != "file-storage" || config.NatsStream != "file-stream" {
		t.Fatalf("Load() = %+v, want dotenv values", *config)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), ".env"))
	if err == nil {
		t.Fatal("Load() error = nil, want an error for a missing file")
	}
}

func writeEnvFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return path
}
