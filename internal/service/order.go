package service

import (
	"context"
	"webook/internal/repository"
	"webook/internal/repository/dao"
)

type OrderService interface {
	List(ctx context.Context) ([]dao.Order, error)
	Add(ctx context.Context, dept dao.Order) error
	Edit(ctx context.Context, dept dao.Order) error
	DeleteByIds(ctx context.Context, ids []int64) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (svc *orderService) List(ctx context.Context) ([]dao.Order, error) {
	return svc.repo.FindAll(ctx)
}

func (svc *orderService) Add(ctx context.Context, dept dao.Order) error {
	return svc.repo.Create(ctx, dept)
}

func (svc *orderService) Edit(ctx context.Context, dept dao.Order) error {
	return svc.repo.Update(ctx, dept)
}

func (svc *orderService) DeleteByIds(ctx context.Context, ids []int64) error {
	return svc.repo.DeleteByIds(ctx, ids)
}
