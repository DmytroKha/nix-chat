package app

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/filesystem"
	"log"
)

type ImageService interface {
	Save(image database.Image, content []byte) (database.Image, error)
	Find(id int64) (database.Image, error)
	FindAll(userId int64) ([]database.Image, error)
	Delete(id int64) error
	Sync(userId int64, objImages, newImages []database.Image) error
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

	allImages, err := s.repo.FindAll(image.UserId)
	if err != nil {
		log.Print(err)
		return database.Image{}, err
	}

	for _, curImage := range allImages {
		if curImage.Id != img.Id {
			err = s.Delete(curImage.Id)
			if err != nil {
				log.Print(err)
				return database.Image{}, err
			}
		}
	}

	return img, nil
}

func (s imageService) Find(id int64) (database.Image, error) {
	return s.repo.Find(id)
}

func (s imageService) FindAll(userId int64) ([]database.Image, error) {
	return s.repo.FindAll(userId)
}

func (s imageService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s imageService) Sync(objId int64, objImages, newImages []database.Image) error {
	imgsNew := make(map[int64]struct{})
	for _, newImage := range newImages {
		isNewImage := true
		imgsNew[newImage.Id] = struct{}{}
		for _, objImage := range objImages {
			if objImage.Id == newImage.Id {
				isNewImage = false
				break
			}
		}
		if isNewImage {
			i, err := s.repo.Find(newImage.Id)
			if err != nil {
				return err
			}
			i.UserId = objId
			_, err = s.repo.Update(i)
			if err != nil {
				return err
			}
		}
	}
	for _, objImage := range objImages {
		if _, exist := imgsNew[objImage.Id]; !exist {
			err := s.repo.Delete(objImage.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
