package repository

import (
	"context"
	"root/internal/team/model"
	"root/pkg/dbs"
)

type ITeamRepository interface {
	Create(ctx context.Context, order *model.Team) error
	FindByName(ctx context.Context, name string) (*model.Team, error)
}

type TeamRepo struct {
	db dbs.IDatabase
}

func NewTeamRepository(db dbs.IDatabase) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team *model.Team) error {
	return r.db.Create(ctx, team)
}

func (r *TeamRepo) FindByName(ctx context.Context, name string) (*model.Team, error) {
	team := new(model.Team)
	opts := []dbs.FindOption{
		dbs.WithQuery(dbs.NewQuery([]string{"team_name = ?"}, name)),
		dbs.WithPreload([]string{"Orders"}),
	}
	// query := dbs.NewQuery([]string{"team_name = ?"}, name)
	if err := r.db.Find(ctx, team, opts...); err != nil {
		return nil, err
	}
	return team, nil
}
