package repository

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(userName, password string) (model.User, error)
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
	GetMessageOwener(messageId int) (int, error)
	UpdateMessage(messageId, userId int, content string) error
	GetMessageById(messageId int) (model.Message, error)
}

type Repository struct {
	Authorization
	Client
	Room
	Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Room:          NewRoomPostgres(db),
		Message:       NewMessagePostgres(db),
		Client:        NewClientPostgres(db),
	}
}
