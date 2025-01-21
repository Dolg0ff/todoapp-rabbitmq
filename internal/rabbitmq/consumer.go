package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/Dolg0ff/todoapp-rabbitmq/internal/metrics"
	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
)

type Consumer struct {
	conn    *Connection
	handler MessageHandler
	logger  *logger.Logger
	metrics *metrics.Metrics
}

type MessageHandler interface {
	Handle(msg []byte) error
}

func NewConsumer(conn *Connection, handler MessageHandler, logger *logger.Logger, metrics *metrics.Metrics) (*Consumer, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}
	if handler == nil {
		return nil, fmt.Errorf("handler is nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	if metrics == nil {
		return nil, fmt.Errorf("metrics is nil")
	}
	if conn.channel == nil {
		return nil, fmt.Errorf("connection channel is nil")
	}

	return &Consumer{
		conn:    conn,
		handler: handler,
		logger:  logger,
		metrics: metrics,
	}, nil
}

func (c *Consumer) Start() error {
	if c == nil {
		return fmt.Errorf("consumer is nil")
	}
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	if c.conn.channel == nil {
		return fmt.Errorf("channel is nil")
	}
	if c.conn.config == nil || c.conn.config.RabbitMQ.Queue == "" {
		return fmt.Errorf("invalid queue configuration")
	}

	c.logger.Info("Starting consumer...")

	// Создаем exchange
	err := c.conn.channel.ExchangeDeclare(
		c.conn.config.RabbitMQ.Exchange, // name: todo_exchange
		"direct",                        // type
		true,                            // durable
		false,                           // auto-deleted
		false,                           // internal
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Создаем очередь
	queue, err := c.conn.channel.QueueDeclare(
		c.conn.config.RabbitMQ.Queue, // name: todo_queue
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Связываем очередь с exchange
	err = c.conn.channel.QueueBind(
		queue.Name,                      // queue name: todo_queue
		c.conn.config.RabbitMQ.Queue,    // routing key: todo_queue
		c.conn.config.RabbitMQ.Exchange, // exchange: todo_exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Start queue size monitoring
	go c.monitorQueueSize()

	msgs, err := c.conn.channel.Consume(
		c.conn.config.RabbitMQ.Queue,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		c.logger.Error("Failed to start consuming", err)
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for msg := range msgs {
			queue, err := c.conn.channel.QueueDeclarePassive(
				c.conn.config.RabbitMQ.Queue,
				true,
				false,
				false,
				false,
				nil,
			)
			if err == nil {
				c.metrics.QueueSize.Set(float64(queue.Messages + 1))
				c.logger.Info("Current queue size", "size", queue.Messages+1)
			}

			c.logger.Info("Received message from queue")

			// Record message size
			c.metrics.MessageSize.Observe(float64(len(msg.Body)))

			// Start processing time measurement
			start := time.Now()

			if err := c.handler.Handle(msg.Body); err != nil {
				c.logger.Error("Failed to handle message", err)
				c.metrics.MessagesFailedTotal.Inc()
				msg.Nack(false, true)
				continue
			}

			queue, err = c.conn.channel.QueueDeclarePassive(
				c.conn.config.RabbitMQ.Queue,
				true,
				false,
				false,
				false,
				nil,
			)
			if err == nil {
				c.metrics.QueueSize.Set(float64(queue.Messages))
				c.logger.Info("Queue size after processing", "size", queue.Messages)
			}

			c.metrics.ProcessingTime.Observe(time.Since(start).Seconds())
			c.metrics.MessagesProcessed.Inc()

			c.logger.Info("Successfully processed message")
			msg.Ack(false)
		}
	}()

	c.logger.Info("Consumer started successfully")

	return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
	c.logger.Info("Stopping consumer...")

	if err := c.conn.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	return nil
}

func (c *Consumer) monitorQueueSize() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		queue, err := c.conn.channel.QueueDeclarePassive(
			c.conn.config.RabbitMQ.Queue, // name
			true,                         // durable
			false,                        // autoDelete
			false,                        // exclusive
			false,                        // noWait
			nil,                          // arguments
		)
		if err != nil {
			c.logger.Error("Failed to inspect queue", err)
			continue
		}
		c.metrics.QueueSize.Set(float64(queue.Messages))
		c.logger.Info("Queue size from monitor", "size", queue.Messages)
	}
}
