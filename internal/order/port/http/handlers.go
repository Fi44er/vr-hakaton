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

func (h *OrderHandler) Update(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	req := new(dto.UpdateOrderReq)
	id := ctx.Params("id")

	if err := ctx.BodyParser(req); err != nil {
		log.Error("Failed to parse body: ", err)
		return response.Error(ctx, 400, err, "Invalid parametrs")
	}

	order, err := h.service.Update(context, id, req)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to update order: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, order)
}

func (h *OrderHandler) Delete(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	id := ctx.Params("id")

	if err := h.service.Delete(context, id); err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to delete order: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, "OK")
}
