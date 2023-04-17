package middlewares

import (
	"fmt"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/resources"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	nethttp "net/http"
)

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
