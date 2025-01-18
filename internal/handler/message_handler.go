package handler

import (
	"time"

	"github.com/Dolg0ff/todoapp-rabbitmq/internal/domain"
	service "github.com/Dolg0ff/todoapp-rabbitmq/internal/service"
	"github.com/Dolg0ff/todoapp-rabbitmq/pkg/logger"
)

type MessageHandler struct {
	service *service.MessageService
	logger  *logger.Logger
}

func NewMessageHandler(service *service.MessageService, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		service: service,
		logger:  logger,
	}
}

func (h *MessageHandler) Handle(msg []byte) error {
	message := &domain.Message{
		Content:   string(msg),
		Timestamp: time.Now(),
	}

	return h.service.ProcessMessage(message)
}
