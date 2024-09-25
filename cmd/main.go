package main

import (
	"root/pkg/config"
	"root/pkg/dbs"

	"root/internal/eventbus"
	orderModel "root/internal/order/model"
	httpServer "root/internal/server/http"
	"root/internal/team/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		log.Fatal("Cannot connnect to database", err)
	}

	if err = db.AutoMigrate(&orderModel.Order{}); err != nil {
		log.Fatal("Database migration fail", err)
	}

  eventBus := eventbus.New()

  teamService := service.NewTeamService()
  teamChanel := make(chan interface{})
  eventBus.Subscribe("order.registred", teamChanel)

  go func() {
    for event := range teamChanel {
      teamService.HandleOrderRegistred(event)
    }
  }()




	validator := validator.New()
	httpSvr := httpServer.NewServer(*validator, db, eventBus)
	if err = httpSvr.Run(); err != nil {
		log.Fatal(err)
	}
}
