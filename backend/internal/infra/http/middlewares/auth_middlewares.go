package middlewares

import (
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetJWTMiddlewares(g *echo.Group, cf config.Configuration) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(cf.JwtSecret),
		ContextKey:    controllers.UserKey,
	}))
}
