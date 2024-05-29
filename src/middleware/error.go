package middleware

import (
	"net/http"
	"tinder-clone/utils/response"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var errCustom *response.Error

	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	switch report.Code {
	case http.StatusNotFound:
		errCustom = response.ErrorWrap(response.CustomError(http.StatusNotFound, "not found"), nil)
	default:
		errCustom = response.ErrorWrap(response.CustomError(http.StatusInternalServerError, "internal server error"), nil)
	}

	response.ErrorResponse(errCustom).Send(c)
}
