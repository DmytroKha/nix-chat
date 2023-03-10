package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(us app.UserService) UserController {
	return UserController{
		userService: us,
	}
}

// FindUser godoc
// @Summary      Show a user
// @Security     ApiKeyAuth
// @Description  get user by name
// @Tags         users
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  resources.UserDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      404  {string}  echo.HTTPError
// @Router       /users/{id} [get]
func (u UserController) Find(ctx echo.Context) error {
	uid := ctx.Param("uid")

	user, err := u.userService.Find(uid)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
	}
	user, err = u.userService.LoadAvatar(user)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var userDto resources.UserDto
	return FormatedResponse(ctx, http.StatusOK, userDto.DatabaseToDto(user))
}

// ChangePass godoc
// @Summary      Change password for user google acc
// @Security     ApiKeyAuth
// @Description  change password for user google acc
// @Tags         users
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "User ID"
// @Param        input   body      requests.ChangePasswordRequest  true  "User body"
// @Success      200  {object}  resources.UserDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /users/{id}/change-pwd [put]
func (u UserController) ChangePassword(ctx echo.Context) error {
	uid := ctx.Param("uid")
	var cpr requests.ChangePasswordRequest
	err := ctx.Bind(&cpr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&cpr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	updatedUser, err := u.userService.ChangePassword(uid, cpr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	updatedUser, err = u.userService.LoadAvatar(updatedUser)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var userDto resources.UserDto
	return FormatedResponse(ctx, http.StatusOK, userDto.DatabaseToDto(updatedUser))
}

// ChangeLogin godoc
// @Summary      Change login for user google acc
// @Security     ApiKeyAuth
// @Description  change login for user google acc
// @Tags         users
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        id   path      string  true  "User ID"
// @Param        input   body      requests.UserRequest  true  "User body"
// @Success      200  {object}  resources.UserDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /users/{id} [put]
func (u UserController) Update(ctx echo.Context) error {
	uid := ctx.Param("uid")

	var usr requests.UserRequest
	err := ctx.Bind(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	updatedUser, err := u.userService.Update(uid, usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	updatedUser, err = u.userService.LoadAvatar(updatedUser)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var userDto resources.UserDto
	return FormatedResponse(ctx, http.StatusOK, userDto.DatabaseToDto(updatedUser))
}
