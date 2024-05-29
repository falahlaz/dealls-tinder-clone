package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"tinder-clone/utils/csvalidator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	var (
		app  = os.Getenv("APP_NAME")
		env  = os.Getenv("ENV")
		name = fmt.Sprintf("%s-%s", app, env)
	)

	e.Use(
		Context,
		echoMiddleware.Recover(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
			AllowHeaders: []string{echo.HeaderAuthorization},
		}),
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: fmt.Sprintf(`{"time":"${time_custom}","remote_ip": ${remote_ip},`+
				`"host":"${host}","method":"${method}","uri":"${uri}","status":${status},`+
				`"error":"${error}","user_agent":"${user_agent}","latency":${latency},"latency_human":"${latency_human}",`+`"server":"main",`+
				`,"name":"%s"}`+"\n", name),
			CustomTimeFormat: time.RFC3339,
			Output:           os.Stdout,
		}),
	)

	e.HTTPErrorHandler = ErrorHandler
	e.Validator = csvalidator.NewCustomValidator()
}
