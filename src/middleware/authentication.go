package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"tinder-clone/src/abstraction"
	"tinder-clone/utils/response"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func parseAuthHeader(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("token not valid")
	}
	return authHeader[7:], nil
}

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header

			authHeader, ok := header["Authorization"]
			if !ok {
				return response.ErrorWrap(response.CustomError(http.StatusUnauthorized, "unauthorized"), nil).Send(c)
			}

			token, err := parseAuthHeader(authHeader[0])
			if err != nil {
				return response.ErrorWrap(response.CustomError(http.StatusUnauthorized, "unauthorized"), err).Send(c)
			}

			authorized, err := IsAuthorized(token, os.Getenv("JWT_SECRET_KEY"))
			if authorized {
				customClaims, err := ExtractData(token, os.Getenv("JWT_SECRET_KEY"))
				if err != nil {
					return response.ErrorWrap(response.CustomError(http.StatusUnauthorized, err.Error()), err).Send(c)
				}

				cc := c.(*abstraction.Context)
				cc.AuthJwt = *customClaims

				return next(c)
			}

			return response.ErrorWrap(response.CustomError(http.StatusUnauthorized, err.Error()), err).Send(c)
		}
	}
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractData(requestToken string, secret string) (*abstraction.JwtCustomClaims, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("invalid Token")
	}

	return &abstraction.JwtCustomClaims{ID: claims["id"].(string), Name: claims["name"].(string)}, nil
}
