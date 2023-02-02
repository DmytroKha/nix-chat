package app

import (
	"errors"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user database.User) (database.User, error)
	Find(id int64) (database.User, error)
	ChangePassword(id int64, cpr requests.ChangePasswordRequest) (database.User, error)
	ChangeLogin(id int64, clr requests.ChangeLoginRequest) (database.User, error)
	GeneratePasswordHash(password string) (string, error)
}

type userService struct {
	userRepo database.UserRepository
}

func NewUserService(r database.UserRepository) UserService {
	return userService{
		userRepo: r,
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
		log.Print(err)
		return database.User{}, err
	}
	return updatedUser, nil
}

func (s userService) ChangeLogin(id int64, clr requests.ChangeLoginRequest) (database.User, error) {
	user, err := s.Find(id)
	if err != nil {
		log.Printf("UserService: %s", err)
		return database.User{}, err
	}
	if clr.NewName == "" {
		return user, errors.New("UserService: user name mismatch")
	}

	user.Name = clr.NewName
	updatedUser, err := s.userRepo.Update(user)
	if err != nil {
		log.Print(err)
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

func (s userService) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
