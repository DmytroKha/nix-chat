package app

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/filesystem"
	"log"
)

type ImageService interface {
	Save(image database.Image, content []byte) (database.Image, error)
	Find(id int64) (database.Image, error)
	Delete(id int64) error
	SaveIntoDB(image database.Image) (database.Image, error)
}

type imageService struct {
	repo    database.ImageRepository
	filesys filesystem.ImageStorageService
}

func NewImageService(r database.ImageRepository, s filesystem.ImageStorageService) ImageService {
	return &imageService{
		repo:    r,
		filesys: s,
	}
}

func (s imageService) Save(image database.Image, content []byte) (database.Image, error) {
	err := s.filesys.SaveImage(image.Name, content)
	if err != nil {
		log.Print(err)
		return database.Image{}, err
	}

	img, err := s.SaveIntoDB(image)
	if err != nil {
		log.Print(err)
		return database.Image{}, err
	}

	return img, nil
}

func (s imageService) Find(id int64) (database.Image, error) {
	return s.repo.Find(id)
}

func (s imageService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s imageService) SaveIntoDB(image database.Image) (database.Image, error) {
	var img database.Image
	img, err := s.repo.Save(image)
	if err != nil {
		log.Print(err)
		return database.Image{}, err
	}
	return img, nil
}
