package service

import (
	"context"
	"root/internal/eventbus"
	"root/internal/order/dto"
	"root/internal/order/model"
	"root/internal/order/repository"
	"root/pkg/mailer"
	"root/pkg/response"
	"root/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type IOrderService interface {
	Register(ctx context.Context, req *dto.RegisterReq) (*model.Order, error)
}

type OrderService struct {
	validator validator.Validate
	repo      repository.IOrderRepository
	eventBus  *eventbus.EventBus
}

func NewOrderService(validator validator.Validate, repo repository.IOrderRepository, eventBus *eventbus.EventBus) *OrderService {
	return &OrderService{
		validator: validator,
		repo:      repo,
		eventBus:  eventBus,
	}
}

func (s *OrderService) Register(ctx context.Context, req *dto.RegisterReq) (*model.Order, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, &response.ErrorResponse{StatusCode: 400, Message: "Invalid params", Err: err}
	}

	order := new(model.Order)
	utils.Copy(order, req)

	resultChan := make(chan eventbus.Result)

	if req.Role == "maintainer" {
		// Create a new team
		if req.Age < 18 {
			return nil, &response.ErrorResponse{StatusCode: 400, Message: "The age of the maintainer person is from 18 years old"}
		}

		s.eventBus.Publish("order.registred", eventbus.OrderRegisteredEvent{TeamName: req.TeamName, ResultChan: resultChan, OrderRole: "maintainer", Context: ctx})
		result := <-resultChan
		if result.Error != nil {
			return nil, result.Error
		}
		if result.Team == nil {
			return nil, &response.ErrorResponse{StatusCode: 400, Message: "Failed to create team"}
		}
		order.TeamID = result.Team.ID
	} else {

		if req.Age < 11 || req.Age > 18 {
			return nil, &response.ErrorResponse{StatusCode: 400, Message: "Unacceptable age"}
		}
		// Check if the team exists
		s.eventBus.Publish("order.registred", eventbus.OrderRegisteredEvent{TeamName: req.TeamName, ResultChan: resultChan, OrderRole: "participant", Context: ctx})
		result := <-resultChan
		if result.Error != nil {
			return nil, result.Error
		}
		if result.Team == nil {
			return nil, &response.ErrorResponse{StatusCode: 404, Message: "Team not found"}
		}
		order.TeamID = result.Team.ID
	}
	close(resultChan)

	existOrder, err := s.repo.FindByEmailOrPhone(ctx, req.Email, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if existOrder.Email != "" {
		return nil, &response.ErrorResponse{StatusCode: 409, Message: "A user with such an email or phone number already exists", Err: err}
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	mailer.Mailer([]string{req.Email}, req.FIO, req.TeamName)

	return order, nil
}
