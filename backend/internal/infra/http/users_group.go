package http

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func UsersGroup(g *echo.Group, userController controllers.UserController, imageController controllers.ImageController) {
	g.GET("/me", userController.Find)
	g.PUT("/change-pwd", userController.ChangePassword)
	g.PUT("/change-name", userController.ChangeName)
	g.PUT("/change_avtr", imageController.AddImage)
}
