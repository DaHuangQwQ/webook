package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) error {
	err := u.dao.Insert(ctx, dao.User{
		Name:     "",
		Age:      0,
		Password: "",
	})
	return err
}
