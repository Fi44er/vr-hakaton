package repository

import "root/pkg/dbs"

type ITeamRepository interface {}

type TeamRepo struct {
  db dbs.IDatabase
}

func NewTeamRepository(db dbs.IDatabase) *TeamRepo {
  return &TeamRepo{db: db}
}

