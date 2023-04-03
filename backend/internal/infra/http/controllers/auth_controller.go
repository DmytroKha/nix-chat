package controllers

import (
	"fmt"
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

	user, token, err := c.authService.Register(usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}
	c.wsServer.AppendUser(&user)

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

	user, token, err := c.authService.Login(usr)
	if err != nil {
		returnErrorResponse(ctx.Response().Writer, http.StatusBadRequest)
		return err
	}
	photo := "../././file_storage/" + user.Image.Name
	//blackList, err := c.userService.GetUserBlackList(user)
	//ctx.Response().Write([]byte(token))
	str := fmt.Sprintf("{\"token\": \"%v\",\"id\": \"%v\",\"photo\": \"%v\"}", token, user.Id, photo)
	ctx.Response().Write([]byte(str))

	return nil

}

func returnErrorResponse(w http.ResponseWriter, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"status\": \"error\"}"))
}
