package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Dolg0ff/todoapp-rabbitmq/internal/config"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/handler"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/rabbitmq"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/service"
	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLogger(cfg.LogLevel)

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	service := service.NewMessageService(logger)
	handler := handler.NewMessageHandler(service, logger)
	consumer, err := rabbitmq.NewConsumer(conn, handler, logger)
	if err != nil {
		logger.Fatal("Failed to create consumer", err)
	}

	if err := consumer.Start(); err != nil {
		logger.Fatal("Failed to start consumer", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logger.Info("Shutting down...")
}
