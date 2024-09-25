package repository

import (
	"context"
	"root/internal/order/model"
	"root/pkg/dbs"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByEmail(ctx context.Context, email string) (*model.Order, error)
}

type OrderRepo struct {
	db dbs.IDatabase
}

func NewOrderRepository(db dbs.IDatabase) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, order *model.Order) error {
	return r.db.Create(ctx, order)
}

func (r *OrderRepo) FindByEmail(ctx context.Context, email string) (*model.Order, error) {
	order := new(model.Order)
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
    return nil, err
	}
  return order, nil
}
