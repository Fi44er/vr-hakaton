package http

import (
	"fmt"
	"root/internal/eventbus"
	orderHTTP "root/internal/order/port/http"
	"root/pkg/config"
	"root/pkg/dbs"
	"root/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Server struct {
	app       *fiber.App
	cfg       *config.Schema
	db        dbs.IDatabase
	validator validator.Validate
  eventBus *eventbus.EventBus
}

func NewServer(validator validator.Validate, db dbs.IDatabase, eventBus *eventbus.EventBus) *Server {
	return &Server{
		app:       fiber.New(),
		cfg:       config.GetConfig(),
		db:        db,
		validator: validator,
    eventBus: eventBus,
	}
}

func (s Server) Run() error {
	// if s.cfg.Environment == config.ProductionEnv {
	//	s.app.Config().Production = true
	// }

	if err := s.MapRoutes(); err != nil {
		log.Fatalf("MapRoutes Error: %v", err)
	}

	// s.app.Get("/swagger/*any", fiberSwagger.WrapHandler(swaggerFiles.Handler))

	s.app.Get("/health", func(c *fiber.Ctx) error {
		response.JSON(c, 200, nil)
		return nil
	})
	// Start http server

	log.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	if err := s.app.Listen(fmt.Sprintf(":%d", s.cfg.HttpPort)); err != nil {
		log.Fatalf("Running HTTP server: %v", err)
	}
	return nil
}

func (s Server) GetApp() *fiber.App {
	return s.app
}

func (s Server) MapRoutes() error {
	v1 := s.app.Group("/api/v1")
	// userHttp.Routes(v1, s.db, s.validator)
  orderHTTP.Routes(v1, s.db, s.validator, s.eventBus)
	return nil
}
