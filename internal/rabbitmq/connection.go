package rabbitmq

import "todoapp-rabbitmq/internal/config"

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.Config
}

func NewConnection(cfg *config.Config) (*Connection, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn:    conn,
		channel: ch,
		config:  cfg,
	}, nil
}

func (c *Connection) Close() {
	c.channel.Close()
	c.conn.Close()
}
