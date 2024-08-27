package article

import (
	"context"
	"webook/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}
