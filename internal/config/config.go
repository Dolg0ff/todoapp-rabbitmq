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

type MetricsConfig struct {
	Port    string `env:"METRICS_PORT"`
	Path    string `env:"METRICS_PATH"`
	Enabled bool   `env:"METRICS_ENABLED"`
}

type Config struct {
	RabbitMQ RabbitMQConfig
	Metrics  MetricsConfig
	LogLevel string `envconfig:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if cfg.RabbitMQ.Queue == "" {
		return fmt.Errorf("RABBITMQ_QUEUE is required")
	}
	if cfg.RabbitMQ.Exchange == "" {
		return fmt.Errorf("RABBITMQ_EXCHANGE is required")
	}

	if cfg.Metrics.Enabled {
		if cfg.Metrics.Port == "" {
			return fmt.Errorf("METRICS_PORT is required when metrics are enabled")
		}
		if cfg.Metrics.Path == "" {
			return fmt.Errorf("METRICS_PATH is required when metrics are enabled")
		}
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s, must be one of: debug, info, warn, error", cfg.LogLevel)
	}

	return nil
}

func (c *Config) String() string {
	maskedURL := maskURL()

	return fmt.Sprintf(
		"Config{RabbitMQ:{URL:%s Queue:%s Exchange:%s} "+
			"Metrics:{Port:%s Path:%s Enabled:%v} "+
			"LogLevel:%s}",
		maskedURL,
		c.RabbitMQ.Queue,
		c.RabbitMQ.Exchange,
		c.Metrics.Port,
		c.Metrics.Path,
		c.Metrics.Enabled,
		c.LogLevel,
	)
}

func maskURL() string {
	return "[MASKED]"
}
