package http

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func AuthGroup(g *echo.Group, authController controllers.AuthController) {
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
}
