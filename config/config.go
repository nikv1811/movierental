package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type MovieAPIHeaders struct {
	RapidAPIHost string `json:"X-RapidAPI-Host"`
	RapidAPIKey  string `json:"X-RapidAPI-Key"`
}

type MovieAPIConfig struct {
	BaseURL string          `json:"base_url"`
	Headers MovieAPIHeaders `json:"headers"`
}

type Config struct {
	Port        string         `json:"port"`
	Environment string         `json:"environment"`
	GinMode     string         `json:"gin_mode"`
	Database    DatabaseConfig `json:"database"`
	MovieAPI    MovieAPIConfig `json:"movie_api"`
}

var AppConfig *Config

func LoadConfig() error {
	file, err := os.Open("config/config.json")
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	cfg := &Config{}
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config JSON: %w", err)
	}

	AppConfig = cfg
	return nil
}
