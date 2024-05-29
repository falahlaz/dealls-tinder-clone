package http

import (
	"fmt"
	"net/http"
	"os"
	"tinder-clone/src/app/auth"
	"tinder-clone/src/app/order"
	"tinder-clone/src/factory"
	"tinder-clone/src/middleware"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, f *factory.Factory) {
	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s", os.Getenv("APP_NAME"))
		return c.String(http.StatusOK, message)
	})

	v1 := e.Group("/api/v1")

	authMiddleware := middleware.Auth()

	auth.NewHandler(f).Route(v1.Group("/auth"))
	order.NewHandler(f).Route(v1.Group("/orders", authMiddleware))
}
