package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	ErrMissingBlueSkyHandle      = errors.New("BLUESKY_HANDLE is required")
	ErrMissingBlueSkyAppPassword = errors.New("BLUESKY_APP_PASSWORD is required")
	ErrMissingServerAPIKey       = errors.New("SERVER_API_KEY is required")
)

type Config struct {
	Bluesky BlueSkyConfig
	Server  ServerConfig
	Log     LogConfig
}

type BlueSkyConfig struct {
	Handle      string
	AppPassword string
}

type ServerConfig struct {
	APIKey string
	Port   int
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Bluesky: BlueSkyConfig{
			Handle:      getEnv("BLUESKY_HANDLE", ""),
			AppPassword: getEnv("BLUESKY_APP_PASSWORD", ""),
		},
		Server: ServerConfig{
			APIKey: getEnv("SERVER_API_KEY", ""),
			Port:   port,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.Bluesky.Handle == "" {
		return ErrMissingBlueSkyHandle
	}
	if c.Bluesky.AppPassword == "" {
		return ErrMissingBlueSkyAppPassword
	}
	if c.Server.APIKey == "" {
		return ErrMissingServerAPIKey
	}
	return nil
}