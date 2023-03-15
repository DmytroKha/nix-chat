package controllers

import (
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	"github.com/labstack/echo/v4"
	"log"
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

func (u UserController) ChangePassword(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := fmt.Errorf("Not authenticated")
		log.Println(err)
		return err
	}
	user := userCtxValue.(domain.User)
	uid := user.GetId()

	var usr requests.ChangePasswordRequest
	err := ctx.Bind(&usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	err = ctx.Validate(&usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusUnprocessableEntity)
		return err
	}

	updatedUser, err := u.userService.ChangePassword(uid, usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte(updatedUser.Password))

	return nil
}

func (u UserController) ChangeName(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := fmt.Errorf("Not authenticated")
		log.Println(err)
		return err
	}
	user := userCtxValue.(domain.User)
	uid := user.GetId()

	var usr requests.UserRequest
	err := ctx.Bind(&usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	err = ctx.Validate(&usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusUnprocessableEntity)
		return err
	}

	if user.GetName() == usr.Name {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	updatedUser, err := u.userService.ChangeName(uid, usr.Name)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte(updatedUser.Password))

	return nil
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
