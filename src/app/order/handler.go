package order

import (
	"net/http"
	"tinder-clone/src/abstraction"
	"tinder-clone/src/dto"
	"tinder-clone/src/factory"
	"tinder-clone/utils/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	*OrderService
}

func (h *handler) Route(e *echo.Group) {
	e.POST("", h.CreateOrder)
	e.POST("/:orderID/payment", h.Payment)
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		OrderService: NewOrderService(f),
	}
}

func (h *handler) CreateOrder(ctx echo.Context) error {
	cc := ctx.(*abstraction.Context)

	payload := &dto.OrderRequestDto{}
	if err := ctx.Bind(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusBadRequest, err.Error())).Send(ctx)
	}

	if err := ctx.Validate(payload); err != nil {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, err.Error())).Send(ctx)
	}

	order, err := h.OrderService.CreateOrder(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(ctx)
	}

	return response.SuccessBuilder(&response.Success{
		Code:     http.StatusCreated,
		Response: order,
	}, order).Send(ctx)
}

func (h *handler) Payment(ctx echo.Context) error {
	cc := ctx.(*abstraction.Context)

	orderID := ctx.Param("orderID")
	if orderID == "" {
		return response.ErrorResponse(response.CustomError(http.StatusUnprocessableEntity, "invalid orderID")).Send(ctx)
	}

	order, err := h.OrderService.Payment(cc, orderID)
	if err != nil {
		return response.ErrorResponse(err).Send(ctx)
	}

	return response.SuccessResponse(order).Send(ctx)
}
