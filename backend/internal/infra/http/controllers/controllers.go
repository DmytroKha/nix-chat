package controllers

import (
	"github.com/labstack/echo/v4"
)

const (
	UserKey = "user"
)

func FormatedResponse(ctx echo.Context, code int, i interface{}) error {
	ct := ctx.Request().Header.Get("Accept")
	if ct == "text/xml" {
		return ctx.XML(code, i)
	} else {
		return ctx.JSON(code, i)
	}
}
