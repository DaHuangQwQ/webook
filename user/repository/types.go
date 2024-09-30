package repository

import (
	"context"
	"webook/user/domain"
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
