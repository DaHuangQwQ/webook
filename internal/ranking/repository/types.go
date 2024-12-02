package repository

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/ranking/domain"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}
