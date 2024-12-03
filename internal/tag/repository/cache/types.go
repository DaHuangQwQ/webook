package cache

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/tag/domain"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = redis.Nil

type TagCache interface {
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	Append(ctx context.Context, uid int64, tags ...domain.Tag) error
	DelTags(ctx context.Context, uid int64) error
}
