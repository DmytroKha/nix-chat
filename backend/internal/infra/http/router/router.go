package router

import (
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/infra/http"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/middlewares"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	"github.com/labstack/echo/v4"
)

func New(userController controllers.UserController,
	authController controllers.AuthController,
	imageController controllers.ImageController,
	wsServer *websocket.WsServer,
	cf config.Configuration) *echo.Echo {

	e := echo.New()
	e.GET("/ws", WsHandlerFunc(wsServer))

	api := e.Group("/api/v1")
	//ws := api.Group("/ws")
	auth := api.Group("/auth")
	users := api.Group("/users")
	images := users.Group("/:userId/image")

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(api)
	//middlewares.SetJWTMiddlewares(ws, cf)
	middlewares.SetJWTMiddlewares(users, cf)
	//middlewares.SetJWTMiddlewares(ws, cf)

	http.MainGroup(e, authController)
	http.AuthGroup(auth, authController)
	http.UsersGroup(users, userController)
	http.ImageGroup(images, imageController)
	//http.ClientGroup(ws, wsServer)

	return e
}

func WsHandlerFunc(wsServer *websocket.WsServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		return websocket.ServeWs(wsServer, c)
	}
}
