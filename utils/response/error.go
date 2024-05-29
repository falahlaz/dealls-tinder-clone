package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
}

type Error struct {
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	ErrorMessage error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code %d", e.Code)
}

func CustomError(httpCode int, message string) *Error {
	return &Error{
		Response: errorResponse{
			StatusCode: fmt.Sprintf("%d", httpCode),
			Message:    message,
		},
		Code: httpCode,
	}
}

func ErrorWrap(errBase *Error, err error) *Error {
	if errBase == nil {
		errBase = CustomError(http.StatusInternalServerError, "internal server error")
	}

	return &Error{
		Response: errorResponse{
			StatusCode: errBase.Response.StatusCode,
			Message:    errBase.Response.Message,
		},
		Code:         errBase.Code,
		ErrorMessage: err,
	}
}

func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorWrap(CustomError(http.StatusInternalServerError, "internal server error"), err)
	}
}

func (e *Error) Send(c echo.Context, request ...interface{}) error {
	return c.JSON(e.Code, e.Response)
}
