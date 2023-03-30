package app

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"log"
)

type BlacklistService interface {
	Save(bl database.Blacklist) (database.Blacklist, error)
	Delete(id int64) error
	Find(userId, foeId int64) (database.Blacklist, error)
	FindAll(userId int64) ([]database.Blacklist, error)
}

type blacklistService struct {
	blRepo database.BlackListRepository
}

func NewBlacklistService(r database.BlackListRepository) BlacklistService {
	return blacklistService{
		blRepo: r,
	}
}

func (s blacklistService) Save(bl database.Blacklist) (database.Blacklist, error) {
	blacklist, err := s.blRepo.Save(bl)
	if err != nil {
		log.Print(err)
		return database.Blacklist{}, err
	}
	return blacklist, nil
}

func (s blacklistService) Delete(id int64) error {
	err := s.blRepo.Delete(id)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (s blacklistService) Find(userId, foeId int64) (database.Blacklist, error) {
	blacklist, err := s.blRepo.Find(userId, foeId)
	if err != nil {
		log.Print(err)
		return database.Blacklist{}, err
	}
	return blacklist, nil
}

func (s blacklistService) FindAll(userId int64) ([]database.Blacklist, error) {
	bls, err := s.blRepo.FindAll(userId)
	if err != nil {
		log.Print(err)
		return []database.Blacklist{}, err
	}
	return bls, nil
}
