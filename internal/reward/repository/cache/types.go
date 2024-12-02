package cache

import (
	"context"
	"github.com/DaHuangQwQ/webook/reward/domain"
)

type RewardCache interface {
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error
}
