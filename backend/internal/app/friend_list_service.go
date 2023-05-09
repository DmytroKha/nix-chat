package app

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"log"
)

//go:generate mockery --dir . --name FriendlistService --output ./mocks
type FriendlistService interface {
	Save(bl database.Friendlist) (database.Friendlist, error)
	Delete(id int64) error
	Find(userId, roomId int64) (database.Friendlist, error)
	FindAll(userId int64) ([]database.Friendlist, error)
}

type friendlistService struct {
	blRepo database.FriendlistRepository
}

func NewFriendlistService(r database.FriendlistRepository) FriendlistService {
	return friendlistService{
		blRepo: r,
	}
}

func (s friendlistService) Save(bl database.Friendlist) (database.Friendlist, error) {
	blacklist, err := s.blRepo.Save(bl)
	if err != nil {
		log.Print("Friendlist: %s", err)
		return database.Friendlist{}, err
	}
	return blacklist, nil
}

func (s friendlistService) Delete(id int64) error {
	err := s.blRepo.Delete(id)
	if err != nil {
		log.Print("Friendlist: %s", err)
		return err
	}
	return nil
}

func (s friendlistService) Find(userId, roomId int64) (database.Friendlist, error) {
	blacklist, err := s.blRepo.Find(userId, roomId)
	if err != nil {
		log.Print("Friendlist: %s", err)
		return database.Friendlist{}, err
	}
	return blacklist, nil
}

func (s friendlistService) FindAll(userId int64) ([]database.Friendlist, error) {
	bls, err := s.blRepo.FindAll(userId)
	if err != nil {
		log.Print("Friendlist: %s", err)
		return []database.Friendlist{}, err
	}
	return bls, nil
}
