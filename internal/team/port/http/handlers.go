package http

import (
	"root/internal/team/dto"
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
	// params := ctx.Queries()
	name := ctx.Query("name")
	// filter := dto.TeamFilterParam{}
	// mapstructure.Decode(params, &filter)

	team, err := h.service.GetWhithPreload(context, name)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to load team: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, team)
}

func (h *TeamHandler) Update(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	req := new(dto.UpdateTeamReq)
	id := ctx.Params("id")

	if err := ctx.BodyParser(req); err != nil {
		log.Error("Failed to parse body: ", err)
		return response.Error(ctx, 400, err, "Invalid parametrs")
	}

	team, err := h.service.Update(context, id, req)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to update team: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, team)
}

func (h *TeamHandler) Delete(ctx *fiber.Ctx) error {
	context := ctx.UserContext()
	id := ctx.Params("id")

	if err := h.service.Delete(context, id); err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to delete order: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, "OK")
}

func (h *TeamHandler) GetAll(ctx *fiber.Ctx) error {
	context := ctx.UserContext()

	teams, err := h.service.GetAll(context)
	if err != nil {
		code, message := response.GetErroField(err)
		log.Error("Failed to getAll teams: ", err)
		return response.Error(ctx, code, err, message)
	}

	return response.JSON(ctx, 200, teams)
}
