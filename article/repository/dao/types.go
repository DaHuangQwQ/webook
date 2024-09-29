package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Upsert(ctx context.Context, article PublishedArticle) error
	Sync(ctx context.Context, article Article) (int64, error)
	SyncStatus(ctx context.Context, articleID int64, authorId int64, status uint8) error
	GetList(ctx context.Context) ([]PublishedArticle, error)
	FindById(ctx context.Context, articleId int64) (PublishedArticle, error)
	GetListByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	ListAll(ctx context.Context, PageNum int, PageSize int) ([]PublishedArticle, error)
	DeleteByIds(ctx context.Context, ids []int64) error
	ListPub(ctx context.Context, start time.Time, pageNum int, pageSize int) ([]PublishedArticle, error)
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&Article{})
}
