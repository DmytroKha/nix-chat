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

var AuthErrInvalidCredentials = errors.New("invalid credentials")

//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	Register(user requests.UserRegistrationRequest) (database.User, string, error)
	Login(user requests.UserLoginRequest) (database.User, string, error)
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

func (s authService) Register(usr requests.UserRegistrationRequest) (database.User, string, error) {
	u, err := usr.ToDatabaseModel()
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	_, err = s.userService.FindByName(u.Name)
	if err == nil {
		log.Printf("AuthService: login error %s", AuthErrInvalidCredentials)
		return database.User{}, "", AuthErrInvalidCredentials
	}

	user, err := s.userService.Save(u)
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	token, err := s.GenerateJwt(user)
	return user, token, err
}

func (s authService) Login(usr requests.UserLoginRequest) (database.User, string, error) {
	user, err := usr.ToDatabaseModel()
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	u, err := s.userService.FindByName(user.Name)
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	u, err = s.userService.LoadAvatar(u)
	if err != nil {
		log.Printf("AuthService: login error %s", err)
		return database.User{}, "", err
	}
	valid := s.checkPasswordHash(user.Password, u.Password)
	if !valid {
		log.Printf("AuthService: login error %s", AuthErrInvalidCredentials)
		return database.User{}, "", AuthErrInvalidCredentials
	}
	token, err := s.GenerateJwt(u)
	return u, token, err
}

func (s authService) GenerateJwt(user database.User) (string, error) {
	claims := resources.JwtClaims{
		ID:    user.Id,
		Name:  user.Name,
		Photo: user.Image.Name,
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
