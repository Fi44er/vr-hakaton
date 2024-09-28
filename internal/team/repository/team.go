package repository

import (
	"context"
	"root/internal/team/model"
	"root/pkg/dbs"
)

type ITeamRepository interface {
	Create(ctx context.Context, order *model.Team) error
	GetById(ctx context.Context, id string) (*model.Team, error)
	FindByName(ctx context.Context, name string) (*model.Team, error)
	FindAll(ctx context.Context) ([]*model.Team, error)
	Update(ctx context.Context, req *model.Team) error
	Delete(ctx context.Context, id string) error
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
		dbs.WithQuery(dbs.NewQuery("team_name = ?", name)),
		dbs.WithPreload([]string{"Orders"}),
	}
	// query := dbs.NewQuery([]string{"team_name = ?"}, name)
	if err := r.db.Find(ctx, team, opts...); err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepo) GetById(ctx context.Context, id string) (*model.Team, error) {
	order := new(model.Team)
	query := dbs.NewQuery("id  = ?", id)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *TeamRepo) FindAll(ctx context.Context) ([]*model.Team, error) {
	teams := make([]*model.Team, 0)
	opts := []dbs.FindOption{
		dbs.WithPreload([]string{"Orders"}),
	}
	if err := r.db.Find(ctx, &teams, opts...); err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *TeamRepo) Update(ctx context.Context, order *model.Team) error {
	return r.db.Update(ctx, order)
}

func (r *TeamRepo) Delete(ctx context.Context, id string) error {
	team := new(model.Team)
	query := dbs.NewQuery("id = ?", id)
	return r.db.Delete(ctx, team, dbs.WithQuery(query))
}
