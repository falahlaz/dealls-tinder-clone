package auth

import (
	"net/http"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/utils/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *AuthService
}

func (h *handler) Route(e *echo.Group) {
	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewAuthService(f),
	}
}

func (h *handler) Login(ctx echo.Context) error {
	cc := ctx.(*abstraction.Context)

	payload := new(dto.LoginRequestDto)
	if err := ctx.Bind(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusBadRequest, err.Error())).Send(ctx)
	}

	if err := ctx.Validate(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, err.Error())).Send(ctx)
	}

	data, err := h.service.Login(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(ctx)
	}

	return response.SuccessResponse(data).Send(ctx)
}

func (h *handler) Register(ctx echo.Context) error {
	cc := ctx.(*abstraction.Context)

	payload := new(dto.RegisterRequestDto)
	if err := ctx.Bind(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusBadRequest, err.Error())).Send(ctx)
	}

	if err := ctx.Validate(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, err.Error())).Send(ctx)
	}

	data, err := h.service.Register(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(ctx)
	}

	return response.SuccessBuilder(&response.Success{
		Code:     http.StatusCreated,
		Response: data,
	}, data).Send(ctx)
}
