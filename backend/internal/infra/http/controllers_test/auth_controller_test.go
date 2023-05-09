package controllers_test

import (
	"bytes"
	"github.com/DmytroKha/nix-chat/internal/app/mocks"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var testUser = database.User{
	Id:   1,
	Name: "testuser",
	Image: database.Image{
		Name: "test.jpg",
	},
}

func TestAuthController_HandleRegister(t *testing.T) {
	e := echo.New()
	e.Validator = requests.NewValidator()
	reqBody := `{"username": "testuser", "password": "123456789", "confirmPassword": "123456789"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	roomService := mocks.NewRoomService(t)
	blacklistService := mocks.NewBlacklistService(t)
	friendlistService := mocks.NewFriendlistService(t)
	authService := mocks.NewAuthService(t)
	userService := mocks.NewUserService(t)
	wsServer := websocket.NewWebsocketServer(roomService, userService, blacklistService, friendlistService)

	controller := controllers.NewAuthController(authService, userService, wsServer)

	authService.On("Register", mock.Anything).Return(testUser, "token123", nil)

	// Test success case
	err := controller.HandleRegister(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "token123", rec.Body.String())

	//authService.AssertCalled(t, "Register", mock.Anything)
	// Test invalid request body
	reqBody = `{123}`
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleRegister(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test unprocessable entity
	reqBody = `{"username": "", "password": "password", "confirmPassword": "password"}`
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleRegister(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Test password and confirm password mismatch
	reqBody = `{"username": "testuser", "password": "password", "confirmPassword": "not_matching_password"}`
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleRegister(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Test error from auth service
	reqBody = `{"username": "userexist", "password": "123456789", "confirmPassword": "123456789"}`
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleRegister(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test returnErrorResponse function
	w := httptest.NewRecorder()
	returnErrorResponse(w, http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"status\": \"error\"}", w.Body.String())
}

func TestAuthController_HandleLogin(t *testing.T) {
	e := echo.New()
	e.Validator = requests.NewValidator()
	reqBody := `{"username": "testuser", "password": "123456789"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	authService := mocks.NewAuthService(t)
	userService := mocks.NewUserService(t)

	controller := controllers.NewAuthController(authService, userService, nil)

	authService.On("Login", mock.Anything).Return(testUser, "token123", nil)

	// Test success case
	err := controller.HandleLogin(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"token\": \"token123\",\"id\": \"1\",\"photo\": \"../../file_storage/test.jpg\"}", rec.Body.String())

	//authService.AssertCalled(t, "Login", mock.Anything)

	// Test invalid request body
	reqBody = `{123}`
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleLogin(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test unprocessable entity
	reqBody = `{"username": "", "password": "password"}`
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleLogin(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Test error from auth service
	reqBody = `{"username": "userabsent", "password": "123456789"}`
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	err = controller.HandleLogin(ctx)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func returnErrorResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte("{\"status\": \"error\"}"))
}
