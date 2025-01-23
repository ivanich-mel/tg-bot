package config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken string
	DB       PostgresConfig
}

type PostgresConfig struct {
	Driver string
	Url    string
}

func Init() (*Config, error) {
	dbConfig, err := initPostgresConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Postgres config: %w", err)
	}

	botToken := os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("missing required environment variable: TG_BOT_TOKEN")
	}

	return &Config{
		BotToken: botToken,
		DB:       *dbConfig,
	}, nil
}

func initPostgresConfig() (*PostgresConfig, error) {
	requiredVars := []string{
		"DB_DRIVER",
		"DB_URL",
	}

	missingVars := []string{}
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			missingVars = append(missingVars, varName)
		}
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	return &PostgresConfig{
		Driver: os.Getenv("DB_DRIVER"),
		Url:    os.Getenv("DB_URL"),
	}, nil
}
