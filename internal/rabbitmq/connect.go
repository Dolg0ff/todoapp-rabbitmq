package rabbitmq

import (
	"fmt"

	"github.com/Dolg0ff/todoapp-rabbitmq/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.Config
}

func NewConnection(cfg *config.Config) (*Connection, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // Закрываем соединение при ошибке создания канала
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Connection{
		conn:    conn,
		channel: ch,
		config:  cfg,
	}, nil
}

func (c *Connection) Close() {
	if c == nil {
		return
	}

	if c.channel != nil {
		c.channel.Close()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}
