package repository

import (
	"context"
	"database/sql"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

var (
	ErrDuplicate    = dao.ErrDuplicate
	ErrUserNotFound = dao.ErrRecordNotFound
)

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao,
		cache,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	err := repo.dao.Insert(ctx, repo.domainToEntity(user))
	return err
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CachedUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
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
	user := repo.entityToDomain(ue)

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

func (repo *CachedUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
	}
}

func (repo *CachedUserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Phone:    user.Phone.String,
	}
}
