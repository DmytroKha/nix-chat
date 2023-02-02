package http

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
)

func ImageGroup(g *echo.Group, imageController controllers.ImageController) {
	g.POST("", imageController.AddImage)
	g.DELETE("/:id", imageController.DeleteImage)
}
