package repository

import (
	"context"
	"github.com/DaHuangQwQ/gutil/slice"
	"webook/search/domain"
	"webook/search/repository/dao"
)

type userRepository struct {
	dao dao.UserDAO
}

func (u *userRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	users, err := u.dao.Search(ctx, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(users, func(idx int, src dao.User) domain.User {
		return domain.User{
			Id:       src.Id,
			Email:    src.Email,
			Nickname: src.Nickname,
			Phone:    src.Phone,
		}
	}), nil
}

func (u *userRepository) InputUser(ctx context.Context, msg domain.User) error {
	return u.dao.InputUser(ctx, dao.User{
		Id:       msg.Id,
		Email:    msg.Email,
		Nickname: msg.Nickname,
		Phone:    msg.Phone,
	})

}

func NewUserRepository(d dao.UserDAO) UserRepository {
	return &userRepository{
		dao: d,
	}
}
