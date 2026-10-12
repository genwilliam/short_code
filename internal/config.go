package internal

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type DatabaseConfig struct {
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	SchemaFile string `yaml:"schema_file"`
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open config %q: %w (copy configs/config.example.yaml to configs/config.yaml and edit it)", path, err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("decode config %q: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return Config{}, fmt.Errorf("decode config %q: multiple YAML documents are not supported", path)
		}
		return Config{}, fmt.Errorf("decode config %q: %w", path, err)
	}

	if config.Server.Address == "" {
		config.Server.Address = ":8080"
	}
	if config.Database.User == "" || config.Database.Host == "" || config.Database.SchemaFile == "" {
		return Config{}, fmt.Errorf("config %q must set database.user, database.host, and database.schema_file", path)
	}
	if config.Database.Port < 1 || config.Database.Port > 65535 {
		return Config{}, fmt.Errorf("config %q database.port must be between 1 and 65535", path)
	}

	if !filepath.IsAbs(config.Database.SchemaFile) {
		config.Database.SchemaFile = filepath.Clean(filepath.Join(filepath.Dir(path), config.Database.SchemaFile))
	}

	return config, nil
}

func (config DatabaseConfig) DSN() string {
	return (&mysql.Config{
		User:      config.User,
		Passwd:    config.Password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(config.Host, fmt.Sprint(config.Port)),
		ParseTime: true,
	}).FormatDSN()
}
