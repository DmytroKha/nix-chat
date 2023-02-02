package app

import (
	"errors"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	Register(user requests.UserRequest) (database.User, string, error)
	Login(user requests.UserRequest) (database.User, string, error)
	GenerateJwt(user database.User) (string, error)
}

type authService struct {
	userService UserService
	config      config.Configuration
}

func NewAuthService(us UserService, cf config.Configuration) AuthService {
	return authService{
		userService: us,
		config:      cf,
	}
}

func (s authService) Register(usr requests.UserRequest) (database.User, string, error) {
	u, err := usr.ToDatabaseModel()
	if err != nil {
		log.Print(err)
		return database.User{}, "", err
	}
	_, err = s.userService.Find(u.Id)
	if err == nil {
		log.Printf("invalid credentials")
		return database.User{}, "", errors.New("invalid credentials")
	}
	user, err := s.userService.Save(u)
	if err != nil {
		log.Print(err)
		return database.User{}, "", err
	}
	token, err := s.GenerateJwt(user)
	return user, token, err
}

func (s authService) Login(usr requests.UserRequest) (database.User, string, error) {
	user, err := usr.ToDatabaseModel()
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	u, err := s.userService.Find(user.Id)
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	valid := s.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return database.User{}, "", errors.New("invalid credentials")
	}
	token, err := s.GenerateJwt(u)
	return u, token, err
}

func (s authService) GenerateJwt(user database.User) (string, error) {
	claims := resources.JwtClaims{
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(s.config.JwtTTL).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s authService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
