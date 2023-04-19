package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
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

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body requests.ChangePasswordRequest true "Data for change password"
// @Success 200 {string} string "User changed password"
// @Failure 400 {string} string "Bad Request"
// @Failure 422 {string} string "Unprocessable Entity"
// @Router /users/password [put]
func (u UserController) ChangePassword(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := AuthErrNotAuth
		log.Println(err)
		return err
	}
	user := userCtxValue.(domain.User)
	id := user.GetId()

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

	updatedUser, err := u.userService.ChangePassword(id, usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte(updatedUser.Password))

	return nil
}

// ChangeName godoc
// @Summary Change user name
// @Description Change user name by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body requests.UserRequest true "Data for change name"
// @Success 200 {string} string "User name after change"
// @Failure 400 {string} string "Bad Request"
// @Failure 422 {string} string "Unprocessable Entity"
// @Router /users/name [put]
func (u UserController) ChangeName(ctx echo.Context) error {

	userCtxValue := ctx.Request().Context().Value("user")
	if userCtxValue == nil {
		err := AuthErrNotAuth
		log.Println(err)
		return err
	}
	user := userCtxValue.(domain.User)
	id := user.GetId()

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

	updatedUser, err := u.userService.ChangeName(id, usr.Name)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	_, err = ctx.Response().Write([]byte(updatedUser.Name))

	return err
}
