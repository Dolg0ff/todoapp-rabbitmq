package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type RabbitMQConfig struct {
	URL      string `envconfig:"RABBITMQ_URL"`
	Queue    string `envconfig:"RABBITMQ_QUEUE"`
	Exchange string `envconfig:"RABBITMQ_EXCHANGE"`
}

type Config struct {
	RabbitMQ RabbitMQConfig
	LogLevel string `envconfig:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if cfg.RabbitMQ.URL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}
	if cfg.RabbitMQ.Queue == "" {
		return nil, fmt.Errorf("RABBITMQ_QUEUE is required")
	}
	if cfg.RabbitMQ.Exchange == "" {
		return nil, fmt.Errorf("RABBITMQ_EXCHANGE is required")
	}
	if cfg.LogLevel == "" {
		return nil, fmt.Errorf("LOG_LEVEL is required")
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	return &cfg, nil
}
