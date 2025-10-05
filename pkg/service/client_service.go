package service

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/repository"
)

type ClientService struct {
	repo repository.Client
}

func NewClientService(repo repository.Client) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) AddClientToRoom(roomId, userId int) error {
	return s.repo.AddClientToRoom(roomId, userId)
}

func (s *ClientService) RemoveClientFromRoom(roomId, userId int) error {
	return s.repo.RemoveClientFromRoom(roomId, userId)
}

func (s *ClientService) GetRoomClients(roomId int) ([]model.User, error) {
	return s.repo.GetRoomClients(roomId)
}
