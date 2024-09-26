package http

import (
	"root/internal/team/service"
	"root/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TeamHandler struct {
	service service.ITeamService
}

func NewOrderHandler(service service.ITeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

func (h *TeamHandler) GetWhithPreload(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	name := ctx.Query("name")

	team, err := h.service.GetWhithPreload(context, name)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to register: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, team)
}
