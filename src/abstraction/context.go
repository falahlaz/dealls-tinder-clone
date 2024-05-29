package abstraction

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type EchoContext interface {
	echo.Context
}

type Context struct {
	EchoContext
	context.Context
	AuthJwt JwtCustomClaims
}

type JwtCustomClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}
