package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tinder-clone/src/factory"
	httpserver "tinder-clone/src/http"
	"tinder-clone/src/middleware"
	"tinder-clone/utils/database"
	"tinder-clone/utils/database/seeder"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Init()

	f := factory.NewFactory()
	e := echo.New()

	middleware.Init(e)
	httpserver.Init(e, f)
	seeder.Init(f)

	url := fmt.Sprintf(`:%s`, os.Getenv("APP_PORT"))
	if err := e.Start(url); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatalf("shutting down the server: %v", err)
	}
}
