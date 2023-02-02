package http

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/requests"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func MainGroup(e *echo.Echo, authController controllers.AuthController) {
	e.Validator = requests.NewValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
