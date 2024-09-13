package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	cache cache.RankingCache
	local *cache.RankingLocalCache
}

func NewCachedRankingRepository(cache cache.RankingCache, local *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{
		cache: cache,
		local: local,
	}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	// 本地缓存 同步 给其他实例
	// 本地 或 redis 缓存预加载 id -> article
	_ = c.local.Set(ctx, arts)
	return c.cache.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := c.local.Get(ctx)
	if err == nil {
		return arts, nil
	}
	arts, err = c.cache.Get(ctx)
	if err == nil {
		return arts, nil
	}
	return arts, err
}
