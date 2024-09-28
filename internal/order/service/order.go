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
	"github.com/gofiber/fiber/v2/log"
)

const ORGANIZER_EMAIL = "ekaterina.dubskaya@yandex.ru"

// const ORGANIZER_EMAIL = "fakeroot94@gmail.com"

type IOrderService interface {
	Register(ctx context.Context, req *dto.RegisterReq) (*model.Order, error)
	Update(ctx context.Context, id string, req *dto.UpdateOrderReq) (*model.Order, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*model.Order, error)
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

	existOrder, err := s.repo.FindByEmailOrPhone(ctx, req.Email, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if existOrder.ID != "" {
		return nil, &response.ErrorResponse{StatusCode: 409, Message: "A user with such an email or phone number already exists", Err: err}
	}

	order := new(model.Order)
	utils.Copy(order, req)

	resultChan := make(chan eventbus.Result)

	if req.Role == "maintainer" {
		// Create a new team
		if req.Age < 18 {
			return nil, &response.ErrorResponse{StatusCode: 400, Message: "The age of the maintainer person is from 18 years old"}
		}

		s.eventBus.Publish("order.registred", eventbus.OrderRegisteredEvent{TeamName: req.TeamName, ResultChan: resultChan, OrderRole: string(req.Role), Context: ctx, Track: string(req.Track)})
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
		s.eventBus.Publish("order.registred", eventbus.OrderRegisteredEvent{TeamName: req.TeamName, ResultChan: resultChan, OrderRole: string(req.Role), Context: ctx})
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

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	mailer.Mailer([]string{req.Email, ORGANIZER_EMAIL}, req.FIO, req.TeamName, order.TeamID, req.Email, req.PhoneNumber, string(req.Role))

	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id string, req *dto.UpdateOrderReq) (*model.Order, error) {
	order, err := s.repo.GetById(ctx, id)
	if err != nil {
		log.Errorf("Update.GetById fail, id: %s, error: %s", id, err)
		return nil, err
	}

	if order.ID == "" {
		return nil, &response.ErrorResponse{StatusCode: 404, Message: "Order not found", Err: err}
	}

	if req.Email != "" || req.PhoneNumber != "" {
		existOrder, err := s.repo.FindByEmailOrPhone(ctx, req.Email, req.PhoneNumber)
		if err != nil {
			return nil, err
		}

		if existOrder.ID != "" {
			return nil, &response.ErrorResponse{StatusCode: 409, Message: "A user with such an email or phone number already exists", Err: err}
		}
	}

	utils.Copy(order, req)

	if err := s.repo.Update(ctx, order); err != nil {
		log.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Errorf("Delete fail, id: %s, error: %s", id, err)
		return err
	}

	return nil
}

func (s *OrderService) GetAll(ctx context.Context) ([]*model.Order, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, &response.ErrorResponse{StatusCode: 404, Message: "Order not found", Err: nil}
	}

	return orders, nil
}
