package router

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/middlewares"
	"github.com/labstack/echo/v4"
)

func New(
	userController controllers.UserController,
	authController controllers.AuthController,
	imageController controllers.ImageController,
) *echo.Echo {

	e := echo.New()

	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	users := api.Group("/users")
	images := users.Group("/:userId/image")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(api)

	http.MainGroup(e, authController)
	http.AuthGroup(auth, authController)
	http.UsersGroup(users, userController)
	http.ImageGroup(images, imageController)

	return e
}
