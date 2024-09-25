package service

import (
	"context"
	"fmt"
	"root/internal/eventbus"
	"root/internal/order/dto"
	"root/internal/order/model"
	"root/internal/order/repository"
	"root/pkg/response"

	"github.com/go-playground/validator/v10"
)

type IOrderService interface {
	Register(ctx context.Context, req *dto.RegisterReq) (*model.Order, error)
}

type OrderService struct {
	validator validator.Validate
	repo      repository.IOrderRepository
  eventBus *eventbus.EventBus
}

func NewOrderService(validator validator.Validate, repo repository.IOrderRepository, eventBus *eventbus.EventBus) *OrderService {
	return &OrderService{
		validator: validator,
		repo:      repo,
    eventBus: eventBus,
	}
}

func (s *OrderService) Register(ctx context.Context, req *dto.RegisterReq) (*model.Order, error) {
	//if err := s.validator.Struct(req); err != nil {
		//return nil, &response.ErrorResponse{StatusCode: 400, Message: "Invalid params", Err: err}
	// }

	order := new(model.Order)
	existOrder, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existOrder.Email != "" {
		return nil, &response.ErrorResponse{StatusCode: 409, Message: "A user with such an email already exists", Err: err}
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

  fmt.Println(req)
  s.eventBus.Publish("order.registred", eventbus.OrderRegisteredEvent{Email: req.Email})

	return order, nil
}
