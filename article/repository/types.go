package repository

import (
	"context"
	"time"
	"webook/article/domain"
	"webook/bff/api"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, articleId int64, authorId int64, status domain.ArticleStatus) error
	GetList(ctx context.Context) ([]domain.Article, error)
	ListById(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, articleId int64) (domain.Article, error)

	Img_Update(ctx context.Context, file []byte, fileType string) (string, error)
	ListAll(ctx context.Context, req api.PageReq) ([]domain.Article, error)
	DeleteByIds(ctx context.Context, ids []int64) error
	ListPub(ctx context.Context, start time.Time, pageNum int, pageSize int) ([]domain.Article, error)
}
