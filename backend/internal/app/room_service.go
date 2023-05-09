package app

import (
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
)

type roomService struct {
	roomRepo database.RoomRepository
}

//go:generate mockery --dir . --name RoomService --output ./mocks
type RoomService interface {
	Save(room domain.Room) (domain.Room, error)
	FindByName(name string) (domain.Room, error)
	FindAll() ([]domain.Room, error)
}

func NewRoomService(roomRepo database.RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (rs *roomService) Save(r domain.Room) (domain.Room, error) {
	var room database.Room

	room.Name = r.GetName()
	room.Private = r.GetPrivate()
	createdRoom, err := rs.roomRepo.Save(room)
	if err != nil {
		return nil, err
	}

	return &createdRoom, nil
}

func (rs *roomService) FindByName(name string) (domain.Room, error) {

	room, err := rs.roomRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (rs *roomService) FindAll() ([]domain.Room, error) {
	rms, err := rs.roomRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var rooms []domain.Room
	for _, rm := range rms {
		room := rm
		rooms = append(rooms, &room)
	}

	return rooms, nil
}
