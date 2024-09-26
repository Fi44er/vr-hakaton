package http

import (
	"root/internal/order/dto"
	"root/internal/order/service"
	"root/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OrderHandler struct {
	service service.IOrderService
}

func NewOrderHandler(service service.IOrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) Register(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	req := new(dto.RegisterReq)
	if err := ctx.BodyParser(req); err != nil {
		log.Error("Failed to parse body: ", err)
		return response.Error(ctx, 400, err, "Invalid parametrs")
	}

	order, err := h.service.Register(context, req)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to register: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, order)
}
