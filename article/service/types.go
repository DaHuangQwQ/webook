package service

import (
	"context"
	"mime/multipart"
	"time"
	"webook/article/domain"
	"webook/bff/api"
)

//go:generate mockgen -source=article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (id int64, err error)
	Img_Update(ctx context.Context, file multipart.File, fileType string) (string, error)
	Withdraw(ctx context.Context, article domain.Article) error
	Publish(ctx context.Context, article domain.Article) (int64, error)
	GetList(ctx context.Context) (list []domain.Article, err error)
	List(ctx context.Context, req api.PageReq) (list []domain.Article, err error)
	ListPub(ctx context.Context, start time.Time, pageNum, pageSize int) ([]domain.Article, error)
	GetPublishedById(ctx context.Context, uid, articleId int64) (domain.Article, error)
	DeleteByIds(ctx context.Context, ids []int64) error
}
