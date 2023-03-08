package controllers

import (
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"

	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthController struct {
	authService app.AuthService
	userService app.UserService
	wsServer    *websocket.WsServer
}

func NewAuthController(as app.AuthService, us app.UserService, ws *websocket.WsServer) AuthController {
	return AuthController{
		authService: as,
		userService: us,
		wsServer:    ws,
	}
}

/*
// NewUser godoc
// @Summary      Create a new user
// @Description  register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        input   body      requests.UserRegistrationRequest  true  "User body"
// @Success      201  {object}  resources.AuthDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /auth/register [post]
func (c AuthController) Register(ctx echo.Context) error {
	var usr requests.UserRegistrationRequest
	err := ctx.Bind(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	user, token, err := c.authService.Register(usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	c.wsServer.AppendUser(user)
	var authDto resources.AuthDto
	return FormatedResponse(ctx, http.StatusCreated, authDto.DatabaseToDto(token, user))

}


// LogInUser godoc
// @Summary      Log in user
// @Description  log in user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Produce      xml
// @Param        input   body      requests.UserRegistrationRequest  true  "User body"
// @Success      200  {object}  resources.AuthDto
// @Failure      400  {string}  echo.HTTPError
// @Failure      422  {string}  echo.HTTPError
// @Failure      500  {string}  echo.HTTPError
// @Router       /auth/login [post]
func (c AuthController) Login(ctx echo.Context) error {
	var usr requests.UserRegistrationRequest
	err := ctx.Bind(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	err = ctx.Validate(&usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusUnprocessableEntity, err)
	}
	user, token, err := c.authService.Login(usr)
	if err != nil {
		return FormatedResponse(ctx, http.StatusBadRequest, err)
	}
	var authDto resources.AuthDto
	return FormatedResponse(ctx, http.StatusOK, authDto.DatabaseToDto(token, user))
}

*/

func (c AuthController) HandleRegister(ctx echo.Context) error {

	var usr requests.UserRegistrationRequest
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

	if usr.Password != usr.ConfirmPassword {
		returnErrorResponse(ctx.Response().Writer, http.StatusUnprocessableEntity)
		return err
	}

	_, token, err := c.authService.Register(usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte(token))

	return nil

}

func (c AuthController) HandleLogin(ctx echo.Context) error {

	var usr requests.UserLoginRequest
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

	_, token, err := c.authService.Login(usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}

	ctx.Response().Write([]byte(token))

	return nil

}

func returnErrorResponse(w http.ResponseWriter, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"status\": \"error\"}"))
}
