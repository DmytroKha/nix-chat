package requests

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
)

type UserRequest struct {
	Name    string `json:"name" validate:"required"`
	ImageId int64  `json:"imageId"`
}

type UserRegistrationRequest struct {
	Password        string `json:"password" validate:"required,alphanum,gte=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,alphanum,gte=6"`
	Name            string `json:"username" validate:"required"`
}

type UserLoginRequest struct {
	Password string `json:"password" validate:"required,alphanum,gte=6"`
	Name     string `json:"username" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,alphanum,gte=6"`
	NewPassword string `json:"newPassword" validate:"required,alphanum,gte=6"`
}

func (r UserRegistrationRequest) ToDatabaseModel() (database.User, error) {
	return database.User{
		Password: r.Password,
		Name:     r.Name,
	}, nil
}

func (r UserLoginRequest) ToDatabaseModel() (database.User, error) {
	return database.User{
		Password: r.Password,
		Name:     r.Name,
	}, nil
}
