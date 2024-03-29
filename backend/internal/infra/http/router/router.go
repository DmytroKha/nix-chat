package router

import (
	"github.com/DmytroKha/nix-chat/internal/infra/http"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/middlewares"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(userController controllers.UserController,
	authController controllers.AuthController,
	imageController controllers.ImageController,
	wsServer *websocket.WsServer) *echo.Echo {

	e := echo.New()
	e.Use(middleware.Logger())
	//logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer logFile.Close()
	//log.SetOutput(logFile)
	//
	//logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Output: logFile,
	//})
	//
	//e.Use(logger)
	//log.Print("test2")

	e.Use(middleware.CORS())
	//e.Use(middleware.Static("../frontend"))
	e.GET("/ws", WsHandlerFunc(wsServer), middlewares.AuthMiddleware)

	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	users := api.Group("/users", middlewares.AuthMiddleware)

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(api)

	http.MainGroup(e, authController)
	http.AuthGroup(auth, authController)
	http.UsersGroup(users, userController, imageController)

	return e
}

func WsHandlerFunc(wsServer *websocket.WsServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		return websocket.ServeWs(wsServer, c)
	}
}
