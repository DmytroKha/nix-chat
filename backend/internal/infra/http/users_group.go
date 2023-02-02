package http

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func UsersGroup(g *echo.Group, userController controllers.UserController) {
	g.GET("/:id", userController.Find)
}
