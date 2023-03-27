package router

import (
	"context"
	"fmt"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/http"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/middlewares"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	nethttp "net/http"
)

func New(userController controllers.UserController,
	authController controllers.AuthController,
	imageController controllers.ImageController,
	wsServer *websocket.WsServer,
	cf config.Configuration) *echo.Echo {

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Static("../frontend"))
	//e.Use(middleware.Static("../backend/file_storage"))
	e.GET("/ws", WsHandlerFunc(wsServer), AuthMiddleware)

	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	users := api.Group("/users", AuthMiddleware)
	//images := users.Group("/change_avtr", AuthMiddleware)

	middlewares.SetMainMiddlewares(e)
	middlewares.SetApiMiddlewares(api)
	//middlewares.SetJWTMiddlewares(users, cf)

	http.MainGroup(e, authController)
	http.AuthGroup(auth, authController)
	http.UsersGroup(users, userController, imageController)
	//http.ImageGroup(images, imageController)

	return e
}

func WsHandlerFunc(wsServer *websocket.WsServer) echo.HandlerFunc {
	return func(c echo.Context) error {
		return websocket.ServeWs(wsServer, c)
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, tok := c.Request().URL.Query()["bearer"]

		if tok && len(token) == 1 {
			user, err := ValidateToken(token[0])
			if err != nil {
				controllers.FormatedResponse(c, nethttp.StatusForbidden, err)
			} else {
				ctx := context.WithValue(c.Request().Context(), "user", user)
				c.SetRequest(c.Request().WithContext(ctx))
			}

		} else {
			controllers.FormatedResponse(c, nethttp.StatusBadRequest, "Please login")
		}
		return next(c)
	}
}

func ValidateToken(tokenString string) (domain.User, error) {
	var conf = config.GetConfiguration()
	token, err := jwt.ParseWithClaims(tokenString, &resources.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(conf.JwtSecret), nil
	})

	if claims, ok := token.Claims.(*resources.JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
