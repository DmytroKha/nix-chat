package middlewares

import (
	"github.com/labstack/echo/v4"
)

func SetMainMiddlewares(e *echo.Echo) {
	e.Use(serverHeader)
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "NIX_chat/v1.0")
		//c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
		//c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		//c.Response().WriteHeader(http.StatusOK)
		return next(c)
	}
}
