package http

import (
	"root/internal/eventbus"
	"root/internal/team/repository"
	"root/internal/team/service"
	"root/pkg/dbs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, db dbs.IDatabase, validator validator.Validate, eventBus *eventbus.EventBus) {
	teamRepo := repository.NewTeamRepository(db)
	teamService := service.NewTeamService(validator, teamRepo)
	teamHandler := NewOrderHandler(teamService)

	teamChanel := make(chan interface{})
	eventBus.Subscribe("order.registred", teamChanel)

	go func() {
		for event := range teamChanel {
			teamService.HandleOrderRegistred(event)
		}
	}()

	team := r.Group("team")
	team.Get("/get-whith-preload", teamHandler.GetWhithPreload)
	team.Get("/get-all", teamHandler.GetAll)
	team.Put("/update/:id", teamHandler.Update)
	team.Delete("/delete/:id", teamHandler.Delete)
}
