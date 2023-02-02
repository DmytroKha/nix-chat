package requests

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
)

type UserRequest struct {
	Password string `json:"password" validate:"required,alphanum,gte=6"`
	Name     string `json:"name" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,alphanum,gte=6"`
	NewPassword string `json:"newPassword" validate:"required,alphanum,gte=6"`
}

type ChangeLoginRequest struct {
	NewName string `json:"newName" validate:"required"`
}

func (r UserRequest) ToDatabaseModel() (database.User, error) {
	return database.User{
		Password: r.Password,
		Name:     r.Name,
	}, nil
}
