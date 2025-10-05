package service

import (
	"github.com/firstproject/talk-together-app/model"
	"github.com/firstproject/talk-together-app/pkg/repository"
)

type RoomService struct {
	repo repository.Room
}

func NewRoomService(repo repository.Room) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(userId int, room model.Room) (int, error) {
	return s.repo.CreateRoom(userId, room)
}

func (s *RoomService) GetAllRooms(userId int) ([]model.Room, error) {
	return s.repo.GetAllRooms(userId)
}

func (s *RoomService) SearchRoomByName(name string) ([]model.Room, error) {
	return s.repo.SearchRoomByName(name)
}

func (s *RoomService) GetRoomById(roomId int) (model.Room, error) {
	return s.repo.GetRoomById(roomId)
}

func (s *RoomService) UpdateRoom(roomId, userId int, input model.UpdateRoomInput) error {
	return s.repo.UpdateRoom(roomId, userId, input)
}

func (s *RoomService) DeleteRoom(userId, roomId int) error {
	return s.repo.DeleteRoom(userId, roomId)
}
