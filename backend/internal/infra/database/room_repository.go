package database

import (
	"github.com/DmytroKha/nix-chat/internal/domain"
	"gorm.io/gorm"
)

const RoomTableName = "rooms"

type Room struct {
	Id      int64 `gorm:"primary_key;auto_increment;not_null"`
	Name    string
	Private bool
}

type roomRepository struct {
	sess *gorm.DB
}

type RoomRepository interface {
	Save(room domain.Room) (domain.Room, error)
	FindByName(name string) (domain.Room, error)
	FindAll() ([]domain.Room, error)
}

func NewRoomRepository(dbSession *gorm.DB) RoomRepository {
	return &roomRepository{
		sess: dbSession,
	}
}

func (room *Room) GetId() int64 {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

func (rr roomRepository) Save(r domain.Room) (domain.Room, error) {
	var room Room

	room.Name = r.GetName()
	room.Private = r.GetPrivate()
	err := rr.sess.Table(RoomTableName).Create(&room).Error
	if err != nil {
		return nil, err
	}

	var rooms []domain.Room
	rooms = append(rooms, &room)
	r = rooms[0]

	return r, nil
}

func (rr roomRepository) FindByName(name string) (domain.Room, error) {
	var r Room
	err := rr.sess.Table(RoomTableName).First(&r, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (rr roomRepository) FindAll() ([]domain.Room, error) {
	var rms []Room
	err := rr.sess.Table(RoomTableName).Where("private = 0").Find(&rms).Error
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
