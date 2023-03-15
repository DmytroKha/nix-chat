package app

import (
	"errors"
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user database.User) (database.User, error)
	Find(uid string) (database.User, error)
	FindByName(name string) (database.User, error)
	ChangePassword(uid string, cpr requests.ChangePasswordRequest) (database.User, error)
	ChangeName(uid, name string) (database.User, error)
	Update(uid string, usr requests.UserRequest) (database.User, error)
	LoadAvatar(user database.User) (database.User, error)
	GeneratePasswordHash(password string) (string, error)
}

type userService struct {
	userRepo     database.UserRepository
	imageService ImageService
}

func NewUserService(r database.UserRepository, is ImageService) UserService {
	return userService{
		userRepo:     r,
		imageService: is,
	}
}

func (s userService) Save(u database.User) (database.User, error) {
	var err error
	u.Password, err = s.GeneratePasswordHash(u.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	user, err := s.userRepo.Save(u)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return user, nil
}

func (s userService) ChangePassword(uid string, cpr requests.ChangePasswordRequest) (database.User, error) {
	user, err := s.Find(uid)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cpr.OldPassword))
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	user.Password, err = s.GeneratePasswordHash(cpr.NewPassword)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) ChangeName(uid, name string) (database.User, error) {
	emptyUser := database.User{}
	user, err := s.Find(uid)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	anotherUser, err := s.FindByName(name)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	if anotherUser != emptyUser {
		err = fmt.Errorf("this name is alreadyused")
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	user.Name = name

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) Update(uid string, usr requests.UserRequest) (database.User, error) {
	var (
		err    error
		image  database.Image
		images []database.Image
	)

	user, err := s.Find(uid)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	if usr.ImageId != 0 {
		var newImages []database.Image
		images, err = s.imageService.FindAll(user.Id)
		if err != nil {
			log.Print(err)
			return database.User{}, err
		}
		if usr.ImageId > 0 {
			image, err = s.imageService.Find(usr.ImageId)
			if err != nil {
				log.Print(err)
				return database.User{}, err
			}
			if user.Id != image.UserId {
				err = errors.New(("this image doesn't belong you. Choose another imageId"))
				return database.User{}, err
			}
			newImages = []database.Image{image}
		} else {
			newImages = []database.Image{}
		}
		err = s.imageService.Sync(user.Id, images, newImages)
		if err != nil {
			log.Print(err)
			return database.User{}, err
		}
	}

	if usr.Name != "" && user.Name != usr.Name {
		user.Name = usr.Name
	}

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) Find(uid string) (database.User, error) {
	user, err := s.userRepo.Find(uid)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return user, err
}

func (s userService) FindByName(name string) (database.User, error) {
	user, err := s.userRepo.FindByName(name)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return user, err
}

func (s userService) LoadAvatar(user database.User) (database.User, error) {
	images, err := s.imageService.FindAll(user.Id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	if len(images) > 0 {
		user.Image = images[0]
	}

	return user, nil
}

func (s userService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
