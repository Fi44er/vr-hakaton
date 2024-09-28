package repository

import (
	"context"
	"root/internal/order/model"
	"root/pkg/dbs"
)

type IOrderRepository interface {
	GetById(ctx context.Context, id string) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
	FindByEmailOrPhone(ctx context.Context, email string, phone string) (*model.Order, error)
	Update(ctx context.Context, req *model.Order) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context) ([]*model.Order, error)
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

func (r *OrderRepo) GetById(ctx context.Context, id string) (*model.Order, error) {
	order := new(model.Order)
	query := dbs.NewQuery("id  = ?", id)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepo) FindByEmailOrPhone(ctx context.Context, email string, phone string) (*model.Order, error) {
	order := new(model.Order)
	query := dbs.NewQuery("email = ?", email)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
		return nil, err
	}

	query = dbs.NewQuery("phone_number  = ?", phone)
	if err := r.db.Find(ctx, order, dbs.WithQuery(query)); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepo) Update(ctx context.Context, order *model.Order) error {
	return r.db.Update(ctx, order)
}

func (r *OrderRepo) Delete(ctx context.Context, id string) error {
	order := new(model.Order)
	query := dbs.NewQuery("id = ?", id)
	return r.db.Delete(ctx, order, dbs.WithQuery(query))
}

func (r *OrderRepo) FindAll(ctx context.Context) ([]*model.Order, error) {
	orders := make([]*model.Order, 0)
	if err := r.db.Find(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
