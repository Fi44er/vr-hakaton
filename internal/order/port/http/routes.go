package http

import (
	"root/internal/eventbus"
	"root/internal/order/repository"
	"root/internal/order/service"
	"root/pkg/dbs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, db dbs.IDatabase, validator validator.Validate, eventBus *eventbus.EventBus) {
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(validator, orderRepo, eventBus)
	orderHandler := NewOrderHandler(orderService)

	order := r.Group("order")
	order.Post("/register", orderHandler.Register)
	order.Put("/update/:id", orderHandler.Update)
	order.Delete("/delete/:id", orderHandler.Delete)
}
