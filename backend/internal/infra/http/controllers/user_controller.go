package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	user, err := u.userService.Find(id)
	if err != nil {
		return FormatedResponse(ctx, http.StatusNotFound, err)
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
// @Router       /users/{id} [put]
func (c UserController) ChangePassword(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}

	var cpr requests.ChangePasswordRequest
	err = ctx.Bind(&cpr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&cpr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	updatedUser, err := c.userService.ChangePassword(id, cpr)
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
// @Param        input   body      requests.ChangeLoginRequest  true  "User body"
// @Success      200  {object}  resources.UserDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /users/{id} [put]
func (c UserController) ChangeLogin(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}

	var clr requests.ChangeLoginRequest
	err = ctx.Bind(&clr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&clr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	updatedUser, err := c.userService.ChangeLogin(id, clr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusInternalServerError, err)
	}
	var userDto resources.UserDto
	return FormatedResponse(ctx, http.StatusOK, userDto.DatabaseToDto(updatedUser))
}
