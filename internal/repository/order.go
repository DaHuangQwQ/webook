package repository

import (
	"context"
	"webook/internal/repository/dao"
)

type OrderRepository interface {
	FindById(ctx context.Context, id int64) (dao.Order, error)
	FindAll(ctx context.Context) ([]dao.Order, error)
	Create(ctx context.Context, order dao.Order) error
	DeleteByIds(ctx context.Context, ids []int64) error
	Update(ctx context.Context, order dao.Order) error
}

type CachedOrderRepository struct {
	dao dao.OrderDao
}

func NewCachedOrderRepository(dao dao.OrderDao) OrderRepository {
	return &CachedOrderRepository{
		dao: dao,
	}
}

func (repo *CachedOrderRepository) FindById(ctx context.Context, id int64) (dao.Order, error) {
	return repo.dao.FindById(ctx, uint64(id))
}

func (repo *CachedOrderRepository) FindAll(ctx context.Context) ([]dao.Order, error) {
	return repo.dao.FindAll(ctx)
}

func (repo *CachedOrderRepository) Create(ctx context.Context, order dao.Order) error {
	return repo.dao.Create(ctx, order)
}

func (repo *CachedOrderRepository) DeleteByIds(ctx context.Context, ids []int64) error {
	return repo.dao.DeleteByIds(ctx, ids)
}

func (repo *CachedOrderRepository) Update(ctx context.Context, order dao.Order) error {
	return repo.dao.Update(ctx, order)
}
