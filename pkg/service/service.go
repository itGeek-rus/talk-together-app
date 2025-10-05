package service

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/kafka"
	"github.com/firstproject/talk-together-app/pkg/redis"
	"github.com/firstproject/talk-together-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(userName, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Room interface {
	CreateRoom(userId int, room model.Room) (int, error)
	GetAllRooms(userId int) ([]model.Room, error)
	SearchRoomByName(name string) ([]model.Room, error)
	GetRoomById(roomId int) (model.Room, error)
	UpdateRoom(roomId, userId int, input model.UpdateRoomInput) error
	DeleteRoom(userId, roomId int) error
}

type Client interface {
	AddClientToRoom(roomId, userId int) error
	RemoveClientFromRoom(roomId, userId int) error
	GetRoomClients(roomId int) ([]model.User, error)
}

type Message interface {
	CreateMessage(roomId, userId int, content string) (int, error)
	GetRoomMessages(roomId int) ([]model.Message, error)
	DeleteMessage(messageId, userId int) error
	UpdateMessage(messageId, userId int, content string) error
}

type Service struct {
	Authorization
	Client
	Room
	Message
	Redis *redis.Client
	Kafka *kafka.Producer
}

func NewService(repos *repository.Repository, redisClient *redis.Client, kafkaProducer *kafka.Producer) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Room:          NewRoomService(repos.Room),
		Message:       NewMessageService(repos.Message, kafkaProducer),
		Client:        NewClientService(repos.Client),
		Redis:         redisClient,
		Kafka:         kafkaProducer,
	}
}
