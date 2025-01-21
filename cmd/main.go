package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dolg0ff/todoapp-rabbitmq/internal/config"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/handler"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/metrics"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/rabbitmq"
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/service"
	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.NewLogger(cfg.LogLevel)

	metrics := metrics.NewMetrics("todoapp", "rabbitmq")

	conn, err := rabbitmq.NewConnection(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	var metricsServer *http.Server
	if cfg.Metrics.Enabled {
		metricsServer = setupMetricsServer(cfg, logger)
	}

	messageService := service.NewMessageService(logger)
	messageHandler := handler.NewMessageHandler(messageService, logger)
	consumer, err := rabbitmq.NewConsumer(conn, messageHandler, logger, metrics)
	if err != nil {
		logger.Fatal("Failed to create consumer", err)
	}

	if err := consumer.Start(); err != nil {
		logger.Fatal("Failed to start consumer", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	logger.Info("Received signal: %v, initiating shutdown...", sig)

	performGracefulShutdown(metricsServer, consumer, logger)
}

func setupMetricsServer(cfg *config.Config, logger *logger.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Metrics.Port),
		Handler: mux,
	}

	go func() {
		logger.Info("Starting metrics server on :%s", cfg.Metrics.Port)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error: %v", err)
		}
	}()

	return metricsServer
}

func performGracefulShutdown(metricsServer *http.Server, consumer *rabbitmq.Consumer, logger *logger.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if metricsServer != nil {
		logger.Info("Shutting down metrics server...")
		if err := metricsServer.Shutdown(ctx); err != nil {
			logger.Error("Metrics server shutdown error: %v", err)
		}
	}

	logger.Info("Stopping consumer...")
	if err := consumer.Stop(ctx); err != nil {
		logger.Error("Consumer shutdown error: %v", err)
	}

	logger.Info("Shutdown completed successfully")
}
