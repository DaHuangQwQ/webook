package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao,
		cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
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

func (repo *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, id)

	switch err {
	case nil:
		// cache 命中
		return u, err
	case cache.ErrUserKeyNotExist:
		// cache 未命中
	default:
		// redis 出错

	}

	ue, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		Id:        ue.Id,
		FirstName: "",
		LastName:  "",
		Password:  "",
		Email:     "",
	}

	go func() {
		err = repo.cache.Set(ctx, user)
		if err != nil {
			// 写入缓存
			// 监控
			//return domain.User{}, err
		}
	}()

	return user, err
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:        0,
		FirstName: "",
		LastName:  "",
		Password:  "",
	}
}
