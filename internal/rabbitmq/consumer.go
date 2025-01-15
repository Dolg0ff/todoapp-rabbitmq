package rabbitmq

import "github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"

type Consumer struct {
	conn    *Connection
	handler MessageHandler
	logger  *logger.Logger
}

type MessageHandler interface {
	Handle(msg []byte) error
}

func NewConsumer(conn *Connection, handler MessageHandler, logger *logger.Logger) *Consumer {
	return &Consumer{
		conn:    conn,
		handler: handler,
		logger:  logger,
	}
}

func (c *Consumer) Start() error {
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
		return err
	}

	go func() {
		for msg := range msgs {
			err := c.handler.Handle(msg.Body)
			if err != nil {
				c.logger.Error("Failed to handle message", err)
				msg.Nack(false, true) // Сообщение вернется в очередь
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}
