package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "configs", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatal(err)
	}
	contents := `server:
  address: ":9090"
database:
  user: test-user
  password: test-password
  host: localhost
  port: 3307
  schema_file: ../schema/init.sql
`
	if err := os.WriteFile(configPath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.Server.Address != ":9090" {
		t.Fatalf("server address = %q, want %q", config.Server.Address, ":9090")
	}
	if config.Database.Password != "test-password" {
		t.Fatalf("database password = %q, want configured password", config.Database.Password)
	}
	wantSchemaPath := filepath.Clean(filepath.Join(filepath.Dir(configPath), "../schema/init.sql"))
	if config.Database.SchemaFile != wantSchemaPath {
		t.Fatalf("schema file = %q, want %q", config.Database.SchemaFile, wantSchemaPath)
	}
}

func TestLoadConfigRejectsUnknownFields(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	contents := `server:
  adress: ":8080"
database:
  user: root
  host: localhost
  port: 3306
  schema_file: schema.sql
`
	if err := os.WriteFile(configPath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadConfig(configPath); err == nil {
		t.Fatal("LoadConfig() error = nil, want unknown-field error")
	}
}
