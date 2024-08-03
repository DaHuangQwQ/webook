package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

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

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Password:  "",
	}
}
