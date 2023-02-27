package database

import (
	"github.com/DmytroKha/nix-chat/internal/domain"
	"gorm.io/gorm"
)

const RoomTableName = "chat_list"

type Room struct {
	Id      int64 `gorm:"primary_key;auto_increment;not_null"`
	Name    string
	Private bool
}

type roomRepository struct {
	sess *gorm.DB
}

func NewRoomRepository(dbSession *gorm.DB) domain.RoomRepository {
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
	err := rr.sess.Table(RoomTableName).Create(&r).Error
	if err != nil {
		return nil, err
	}
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
