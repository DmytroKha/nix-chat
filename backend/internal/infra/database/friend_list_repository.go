package database

import (
	"gorm.io/gorm"
)

const FriendlistTableName = "friend_list"

type Friendlist struct {
	Id       int64 `gorm:"primary_key;auto_increment;not_null"`
	UserId   int64 `json:"userId"`
	FriendId int64 `json:"friendId"`
	RoomId   int64 `json:"roomId"`
}

type friendlistRepository struct {
	sess *gorm.DB
}

type FriendlistRepository interface {
	Save(bl Friendlist) (Friendlist, error)
	Delete(id int64) error
	Find(userId, roomId int64) (Friendlist, error)
	FindAll(userId int64) ([]Friendlist, error)
}

func NewFriendlistRepository(dbSession *gorm.DB) FriendlistRepository {
	return &friendlistRepository{
		sess: dbSession,
	}
}

func (r friendlistRepository) Save(bl Friendlist) (Friendlist, error) {
	err := r.sess.Table(FriendlistTableName).Create(&bl).Error
	if err != nil {
		return Friendlist{}, err
	}
	return bl, nil
}

func (r friendlistRepository) Delete(id int64) error {
	err := r.sess.Table(FriendlistTableName).Where("id = ?", id).Delete(Friendlist{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r friendlistRepository) Find(userId, roomId int64) (Friendlist, error) {
	var bl Friendlist
	err := r.sess.Table(FriendlistTableName).First(&bl, "user_id = ? AND room_id = ?", userId, roomId).Error
	if err != nil {
		return Friendlist{}, err
	}
	return bl, nil
}

func (r friendlistRepository) FindAll(userId int64) ([]Friendlist, error) {
	var bls []Friendlist
	err := r.sess.Table(FriendlistTableName).Find(&bls, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return bls, nil
}
