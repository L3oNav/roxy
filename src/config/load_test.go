package config

import (
	"os"
	"testing"
)

// [server]
// uri = "/api"
// logfile = "logs/access.log"
// loglevel = "debug"
// [[forward]]
// algorithm = "WRR"
// backends = [
//     { address = "127.0.0.1:8080", weight = 1 },
//     { address = "127.0.0.1:8081", weight = 3 },
//     { address = "127.0.0.1:8082", weight = 2 },
// ]

func TestLoadConfig(t *testing.T) {
	var err error
	tmpfile, err := os.Create("config.toml")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	configContent := `
		[server]
		name = "roxy"
		logfile = "logs/access.log"
		loglevel = "debug"
		listen = ["127.0.0.1:3312"]
		max_connections = 1024

		[[match]]
		uri = "/"
		serve = "/static"


		[[match]]
		uri = "/api"
		algorithm = "WRR"
		forward = [
			{ address = "127.0.0.1:8080", weight = 1 },
			{ address = "127.0.0.1:8081", weight = 3 },
			{ address = "127.0.0.1:8082", weight = 2 },
		]
	`
	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	config := NewConfig()
	config, err = config.Load(tmpfile.Name())

	if err != nil {
		t.Errorf("Load() error = %v, wantErr %v", err, false)
	}

	if config.Server.NAME != "roxy" {
		t.Errorf("Load() got = %v, want %v", config.Server.NAME, "roxy")
	}

	if len(config.Pattern) != 2 {
		t.Errorf("Load() got = %v, want %v", len(config.Pattern), 2)
	}

}
