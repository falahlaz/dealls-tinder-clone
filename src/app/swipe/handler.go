package swipe

import (
	"net/http"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/utils/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	*SwipeService
}

func (h *handler) Route(e *echo.Group) {
	e.POST("/potential-matches", h.FindPotentialMatches)
	e.POST("/matches", h.Match)
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		SwipeService: NewSwipeService(f),
	}
}

func (h *handler) FindPotentialMatches(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := &dto.PotentialMatchRequestDto{}
	if err := c.Bind(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusBadRequest, err.Error())).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, err.Error())).Send(c)
	}

	user, err := h.SwipeService.FindPotentialMatches(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(user).Send(c)
}

func (h *handler) Match(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := &dto.MatchRequestDto{}
	if err := c.Bind(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusBadRequest, err.Error())).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, err.Error())).Send(c)
	}

	err := h.SwipeService.Match(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(nil).Send(c)
}
