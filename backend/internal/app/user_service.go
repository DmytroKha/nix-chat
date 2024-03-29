package app

import (
	"errors"
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var UserErrChooseAnotherImage = errors.New(("this image doesn't belong you. Choose another imageId"))

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user database.User) (database.User, error)
	Find(id int64) (database.User, error)
	FindByName(name string) (database.User, error)
	ChangePassword(id int64, cpr requests.ChangePasswordRequest) (database.User, error)
	ChangeName(id int64, name string) (database.User, error)
	Update(id int64, usr requests.UserRequest) (database.User, error)
	LoadAvatar(user database.User) (database.User, error)
	GeneratePasswordHash(password string) (string, error)
	GetUserBlackList(user domain.User) ([]domain.User, error)
	GetUserFriends(user domain.User) ([]domain.User, error)
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
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return user, nil
}

func (s userService) ChangePassword(id int64, cpr requests.ChangePasswordRequest) (database.User, error) {
	user, err := s.Find(id)
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
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) ChangeName(id int64, name string) (database.User, error) {
	emptyUser := database.User{}
	user, err := s.Find(id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	anotherUser, _ := s.FindByName(name)

	if anotherUser != emptyUser {
		err = fmt.Errorf("this name is already used")
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	user.Name = name

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) Update(id int64, usr requests.UserRequest) (database.User, error) {
	var (
		err    error
		image  database.Image
		images []database.Image
	)

	user, err := s.Find(id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}

	if usr.ImageId != 0 {
		var newImages []database.Image
		images, err = s.imageService.FindAll(user.Id)
		if err != nil {
			log.Printf("UserService: %s", err)
			return database.User{}, err
		}
		if usr.ImageId > 0 {
			image, err = s.imageService.Find(usr.ImageId)
			if err != nil {
				log.Printf("UserService: %s", err)
				return database.User{}, err
			}
			if user.Id != image.UserId {
				err = UserErrChooseAnotherImage
				log.Printf("UserService: %s", err)
				return database.User{}, err
			}
			newImages = []database.Image{image}
		} else {
			newImages = []database.Image{}
		}
		err = s.imageService.Sync(user.Id, images, newImages)
		if err != nil {
			log.Printf("UserService: %s", err)
			return database.User{}, err
		}
	}

	if usr.Name != "" && user.Name != usr.Name {
		user.Name = usr.Name
	}

	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) Find(id int64) (database.User, error) {
	user, err := s.userRepo.Find(id)
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

func (s userService) GetUserBlackList(user domain.User) ([]domain.User, error) {
	users, err := s.userRepo.GetUserBlackList(user.GetId())
	if err != nil {
		log.Printf("UserService: %s", err)
		return []domain.User{}, err
	}
	var blackList []domain.User
	for _, u := range users {
		newUser := u
		blackList = append(blackList, &newUser)
	}
	return blackList, err
}

func (s userService) GetUserFriends(user domain.User) ([]domain.User, error) {
	users, err := s.userRepo.GetUserFriends(user.GetId())
	if err != nil {
		log.Printf("UserService: %s", err)
		return []domain.User{}, err
	}
	var friends []domain.User
	for _, u := range users {
		newUser := u
		friends = append(friends, &newUser)
	}
	return friends, err
}

func (s userService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
