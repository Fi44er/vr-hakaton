package service

import (
	"context"
	"fmt"
	"root/internal/eventbus"
	"root/internal/team/dto"
	"root/internal/team/model"
	"root/internal/team/repository"
	"root/pkg/response"
	"root/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

type ITeamService interface {
	HandleOrderRegistred(event interface{})
	GetWhithPreload(ctx context.Context, name string) (*model.Team, error)
	GetAll(ctx context.Context) ([]*model.Team, error)
	Update(ctx context.Context, id string, req *dto.UpdateTeamReq) (*model.Team, error)
	Delete(ctx context.Context, id string) error
}

type TeamService struct {
	validator validator.Validate
	repo      repository.ITeamRepository
}

func NewTeamService(validator validator.Validate, repo repository.ITeamRepository) *TeamService {
	return &TeamService{
		validator: validator,
		repo:      repo,
	}
}

func (s *TeamService) HandleOrderRegistred(event interface{}) {
	// Попытка приведения типа event к типу OrderRegisteredEvent
	if e, ok := event.(eventbus.OrderRegisteredEvent); ok {
		// Проверяем, существует ли уже команда с таким названием
		existingTeam, err := s.repo.FindByName(e.Context, e.TeamName)
		if err != nil {
			log.Error("Error checking existing team", err)
			e.ResultChan <- eventbus.Result{Team: nil, Error: err}
			return
		}

		if e.OrderRole == "maintainer" {
			if existingTeam.ID != "" {
				message := fmt.Sprintf("team with name %s already exists", e.TeamName)
				e.ResultChan <- eventbus.Result{Team: nil, Error: &response.ErrorResponse{StatusCode: 409, Message: message, Err: nil}}
				return
			}

			newTeam := &model.Team{
				TeamName: e.TeamName,
				Track:    model.Track(e.Track),
			}

			// Сохранение команды в базе данных
			if err := s.repo.Create(e.Context, newTeam); err != nil {
				log.Error("Failed to create team", err)
				e.ResultChan <- eventbus.Result{Team: nil, Error: err}
				return
			}
			e.ResultChan <- eventbus.Result{Team: newTeam, Error: nil}
		} else {
			if existingTeam.ID == "" {
				message := fmt.Sprintf("team with name %s not found", e.TeamName)
				e.ResultChan <- eventbus.Result{Team: nil, Error: &response.ErrorResponse{StatusCode: 404, Message: message, Err: nil}}
				return
			}

			team, err := s.GetWhithPreload(e.Context, e.TeamName)
			if err != nil {
				e.ResultChan <- eventbus.Result{Team: nil, Error: err}
				return
			}

			if len(team.Orders) == 4 {
				message := fmt.Sprintf("the team %s is staffed", e.TeamName)
				e.ResultChan <- eventbus.Result{Team: nil, Error: &response.ErrorResponse{StatusCode: 409, Message: message, Err: nil}}
				return
			}

			captains := false
			for _, order := range team.Orders {
				if order.Role == "captain" {
					captains = true
					break
				}
			}

			if e.OrderRole == "captain" && captains {
				message := "the captain already exists"
				e.ResultChan <- eventbus.Result{Team: nil, Error: &response.ErrorResponse{StatusCode: 409, Message: message, Err: nil}}
				return
			}

			e.ResultChan <- eventbus.Result{Team: existingTeam, Error: nil}
		}
	}
}

func (s *TeamService) GetWhithPreload(ctx context.Context, name string) (*model.Team, error) {
	team, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if team.ID == "" {
		return nil, &response.ErrorResponse{StatusCode: 404, Message: "Team not found", Err: nil}
	}
	return team, nil
}

func (s *TeamService) GetAll(ctx context.Context) ([]*model.Team, error) {
	teams, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(teams) == 0 {
		return nil, &response.ErrorResponse{StatusCode: 404, Message: "Team not found", Err: nil}
	}

	return teams, nil
}

func (s TeamService) Update(ctx context.Context, id string, req *dto.UpdateTeamReq) (*model.Team, error) {
	team, err := s.repo.GetById(ctx, id)
	if err != nil {
		log.Errorf("Update.GetById fail, id: %s, error: %s", id, err)
		return nil, err
	}

	if team.ID == "" {
		return nil, &response.ErrorResponse{StatusCode: 404, Message: "Team not found", Err: err}
	}

	utils.Copy(team, req)

	if err := s.repo.Update(ctx, team); err != nil {
		log.Errorf("Update fail, id: %s, error: %s", id, err)
		return nil, err
	}

	return team, nil
}

func (s TeamService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Errorf("Delete fail, id: %s, error: %s", id, err)
		return err
	}

	return nil
}
