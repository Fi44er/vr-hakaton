package service

import (
	"root/internal/eventbus"

	"github.com/gofiber/fiber/v2/log"
)

type ITeamService interface {
  HandleOrderRegistred(event interface{})
}

type TeamService struct {
  //validator validator.Validate
  //repo repository.ITeamRepository
}
// validator validator.Validate, repo repository.ITeamRepository
func NewTeamService() *TeamService {
  return &TeamService{
    // validator: validator,
    // repo: repo,
  }
}

func (s *TeamService) HandleOrderRegistred(event interface{}) {
  if e, ok := event.(eventbus.OrderRegisteredEvent); ok {
    log.Infof("Creating team for order with email: %s", e.Email)
  }
}
