package service

import (
	"github.com/Dolg0ff/todoapp-rabbitmq/internal/domain"
	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
)

type MessageService struct {
	logger *logger.Logger
}

func NewMessageService(logger *logger.Logger) *MessageService {
	return &MessageService{
		logger: logger,
	}
}

func (s *MessageService) ProcessMessage(msg *domain.Message) error {
	// Бизнес-логика обработки сообщения
	s.logger.Info("Processing message", "content", msg.Content)
	return nil
}
