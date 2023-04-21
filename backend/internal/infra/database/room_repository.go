package database

import (
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
	Save(room Room) (Room, error)
	FindByName(name string) (Room, error)
	FindAll() ([]Room, error)
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

func (rr roomRepository) Save(r Room) (Room, error) {
	err := rr.sess.Table(RoomTableName).Create(&r).Error
	if err != nil {
		return Room{}, err
	}

	return r, nil
}

func (rr roomRepository) FindByName(name string) (Room, error) {
	var r Room
	err := rr.sess.Table(RoomTableName).First(&r, "name = ?", name).Error
	if err != nil {
		return Room{}, err
	}
	return r, nil
}

func (rr roomRepository) FindAll() ([]Room, error) {
	var rms []Room
	err := rr.sess.Table(RoomTableName).Where("private = 0").Find(&rms).Error
	if err != nil {
		return nil, err
	}

	return rms, nil
}
