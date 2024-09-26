package repository

import (
	"context"
	"root/internal/order/model"
	"root/pkg/dbs"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByEmailOrPhone(ctx context.Context, email string, phone string) (*model.Order, error)
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

func (r *OrderRepo) FindByEmailOrPhone(ctx context.Context, email string, phone string) (*model.Order, error) {
	order := new(model.Order)
	query := dbs.NewQuery([]string{"email = ?", "OR phone_number = ?"}, email, phone)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
		return nil, err
	}
	return order, nil
}
