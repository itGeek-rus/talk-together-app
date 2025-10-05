package service

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/kafka"
	"github.com/firstproject/talk-together-app/pkg/middleware/monitoring"
	"github.com/firstproject/talk-together-app/pkg/repository"
	"time"
)

type MessageService struct {
	repo  repository.Message
	kafka *kafka.Producer
}

func NewMessageService(repo repository.Message, kafkaProducer *kafka.Producer) *MessageService {
	service := &MessageService{
		repo:  repo,
		kafka: kafkaProducer,
	}

	return service
}

func (s *MessageService) CreateMessage(roomId, userId int, content string) (int, error) {
	id, err := s.repo.CreateMessage(roomId, userId, content)
	if err != nil {
		return 0, err
	}

	message := model.Message{
		Id:        id,
		Room:      roomId,
		User:      userId,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if s.kafka != nil {
		err := s.kafka.SendMessage("messages", message)
		if err != nil {
			monitoring.IncrementKafkaMessagesSent("messages_error")
			return 0, nil
		}
		monitoring.IncrementKafkaMessagesSent("messages_success")
	}

	return message.Id, nil
}

func (s *MessageService) GetRoomMessages(roomId int) ([]model.Message, error) {
	return s.repo.GetRoomMessages(roomId)
}

func (s *MessageService) DeleteMessage(messageId, userId int) error {
	return s.repo.DeleteMessage(messageId, userId)
}

func (s *MessageService) UpdateMessage(messageId, userId int, content string) error {
	return s.repo.UpdateMessage(messageId, userId, content)
}
