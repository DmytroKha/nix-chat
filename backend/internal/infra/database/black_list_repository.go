package database

import (
	"gorm.io/gorm"
)

const BlacklistTableName = "black_list"

type Blacklist struct {
	Id     int64 `gorm:"primary_key;auto_increment;not_null"`
	UserId int64 `json:"userId"`
	FoeId  int64 `json:"foeId"`
	RoomId int64 `json:"roomId"`
}

type blacklistRepository struct {
	sess *gorm.DB
}

type BlacklistRepository interface {
	Save(bl Blacklist) (Blacklist, error)
	Delete(id int64) error
	Find(userId, roomId int64) (Blacklist, error)
	FindAll(userId int64) ([]Blacklist, error)
	//GetUserBlackList(user domain.User) ([]domain.User, error)
}

func NewBlacklistRepository(dbSession *gorm.DB) BlacklistRepository {
	return &blacklistRepository{
		sess: dbSession,
	}
}

func (r blacklistRepository) Save(bl Blacklist) (Blacklist, error) {
	err := r.sess.Table(BlacklistTableName).Create(&bl).Error
	if err != nil {
		return Blacklist{}, err
	}
	return bl, nil
}

func (r blacklistRepository) Delete(id int64) error {
	err := r.sess.Table(BlacklistTableName).Where("id = ?", id).Delete(Blacklist{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *blacklistRepository) Find(userId, roomId int64) (Blacklist, error) {
	var bl Blacklist
	err := r.sess.Table(BlacklistTableName).First(&bl, "user_id = ? AND room_id = ?", userId, roomId).Error
	if err != nil {
		return Blacklist{}, err
	}
	return bl, nil
}

func (r blacklistRepository) FindAll(userId int64) ([]Blacklist, error) {
	var bls []Blacklist
	err := r.sess.Table(BlacklistTableName).Find(&bls, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return bls, nil
}
