package controllers

import (
	"errors"
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var AuthErrNotAuth = errors.New("not authenticated")

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

// HandleRegister godoc
// @Summary Register a new user.
// @Description Register a new user with the provided details.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param userRegistrationRequest body requests.UserRegistrationRequest true "User registration details."
// @Success 200 {string} string "User registered successfully."
// @Failure 400 {string} string "Invalid request format or missing required fields."
// @Failure 422 {string} string "Validation errors occurred."
// @Failure 500 {string} string "Internal server error."
// @Router /auth/register [post]
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

// HandleLogin godoc
// @Summary Handle user login
// @Description Log in a user with their credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.UserLoginRequest true "User login details"
// @Success 200 {string} string "Token, user ID, and user photo"
// @Failure 400 {string} string "Bad Request"
// @Failure 422 {string} string "Unprocessable Entity"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
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
	//photo := "../././file_storage/" + user.Image.Name
	photo := "../../file_storage/" + user.Image.Name

	str := fmt.Sprintf("{\"token\": \"%v\",\"id\": \"%v\",\"photo\": \"%v\"}", token, user.Id, photo)
	ctx.Response().Write([]byte(str))

	return nil

}

func returnErrorResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"status\": \"error\"}"))
}
