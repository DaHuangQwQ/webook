package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/repository/dao/oss"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, OpenId string) (domain.User, error)
	UpdateByID(ctx context.Context, user domain.User) error
	AvatarUpdate(ctx context.Context, id int64, file []byte, fileType string) (string, error)
	GetAvatar(ctx context.Context, id int64) (string, error)
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
	oss   oss.Client
}

func (repo *CachedUserRepository) GetAvatar(ctx context.Context, id int64) (string, error) {
	user, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return "", err
	}
	return user.AvatarUrl, nil
}

var (
	ErrDuplicate    = dao.ErrDuplicate
	ErrUserNotFound = dao.ErrRecordNotFound
)

func NewUserRepository(dao dao.UserDao, cache cache.UserCache, oss oss.Client) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
		oss:   oss,
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

	return user, nil
}

func (repo *CachedUserRepository) FindByWechat(ctx context.Context, OpenId string) (domain.User, error) {
	u, err := repo.dao.FindByWechat(ctx, OpenId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *CachedUserRepository) UpdateByID(ctx context.Context, user domain.User) error {
	return repo.dao.Update(ctx, repo.domainToEntity(user))
}

func (repo *CachedUserRepository) AvatarUpdate(ctx context.Context, id int64, file []byte, fileType string) (string, error) {
	type Config struct {
		ENDP string `yaml:"ENDP"`
	}
	var config Config
	err := viper.UnmarshalKey("OSS", &config)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("avatar/%d.%s", id, fileType)
	ossAdress := "https://ceit." + config.ENDP + "/" + fileName
	err = repo.dao.Update(ctx, dao.User{
		Id:        id,
		AvatarUrl: ossAdress,
	})
	if err != nil {
		return "", err
	}

	return ossAdress, repo.oss.UploadFile(ctx, fileName, file)
}

func (repo *CachedUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		WechatOpenId: sql.NullString{
			String: user.WechatInfo.OpenId,
			Valid:  user.WechatInfo.OpenId != "",
		},
		WechatUnionId: sql.NullString{
			String: user.WechatInfo.UnionId,
			Valid:  user.WechatInfo.UnionId != "",
		},
		Grade:  user.Grade,
		Gender: user.Gender,
	}
}

func (repo *CachedUserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Phone:    user.Phone.String,
		Nickname: user.Nickname,
		Grade:    user.Grade,
		Gender:   user.Gender,
		CTime:    time.UnixMilli(user.CTime),

		WechatInfo: domain.WechatInfo{
			OpenId:  user.WechatOpenId.String,
			UnionId: user.WechatUnionId.String,
		},
	}
}
