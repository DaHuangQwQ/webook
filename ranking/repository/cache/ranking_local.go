package cache

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
	"webook/ranking/domain"
)

type RankingLocalCache struct {
	topN *atomicx.Value[[]domain.Article]
	ddl  *atomicx.Value[time.Time]
	exp  time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		topN: atomicx.NewValue[[]domain.Article](),
		ddl:  atomicx.NewValueOf[time.Time](time.Now()),
		exp:  time.Minute * 10,
	}
}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.topN.Store(arts)
	r.ddl.Store(time.Now().Add(r.exp))
	return nil
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := r.ddl.Load()
	arts := r.topN.Load()
	if ddl.Before(time.Now()) || len(arts) == 0 {
		return nil, fmt.Errorf("local cache 已经过期 或 没命中")
	}
	return arts, nil
}
