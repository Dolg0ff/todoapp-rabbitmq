package rabbitmq

import (
	"fmt"

	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
)

type Consumer struct {
	conn    *Connection
	handler MessageHandler
	logger  *logger.Logger
}

type MessageHandler interface {
	Handle(msg []byte) error
}

func NewConsumer(conn *Connection, handler MessageHandler, logger *logger.Logger) (*Consumer, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}
	if handler == nil {
		return nil, fmt.Errorf("handler is nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	if conn.channel == nil {
		return nil, fmt.Errorf("connection channel is nil")
	}

	return &Consumer{
		conn:    conn,
		handler: handler,
		logger:  logger,
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
			c.logger.Info("Received message from queue")
			if err := c.handler.Handle(msg.Body); err != nil {
				c.logger.Error("Failed to handle message", err)
				msg.Nack(false, true) // Сообщение вернется в очередь
				continue
			}
			c.logger.Info("Successfully processed message")
			msg.Ack(false)
		}
	}()

	c.logger.Info("Consumer started successfully")

	return nil
}
