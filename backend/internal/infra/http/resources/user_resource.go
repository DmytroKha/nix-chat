package resources

import (
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	jwt "github.com/dgrijalva/jwt-go"
)

type UserDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UsersDto struct {
	Items []UserDto `json:"items"`
	Total uint64    `json:"total"`
	Pages uint64    `json:"pages"`
}

type AuthDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type GoogleUrlDto struct {
	Url string `json:"url"`
}

func (d UserDto) DatabaseToDto(user database.User) UserDto {
	return UserDto{
		Id:   user.Id,
		Name: user.Name,
	}
}

func (d AuthDto) DatabaseToDto(token string, user database.User) AuthDto {
	var userDto UserDto
	a := AuthDto{
		Token: token,
		User:  userDto.DatabaseToDto(user),
	}
	return a
}

func (d GoogleUrlDto) DatabaseToDto(url string) GoogleUrlDto {
	return GoogleUrlDto{
		Url: url,
	}
}

func (d UserDto) DatabaseToDtoCollection(users database.Users) UsersDto {
	result := make([]UserDto, len(users.Items))
	for i := range users.Items {
		result[i] = d.DatabaseToDto(users.Items[i])
	}
	return UsersDto{Items: result, Pages: users.Pages, Total: users.Total}
}