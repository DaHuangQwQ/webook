package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewGormArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

func (dao *GormArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CTime = now
	article.UTime = now

	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (dao *GormArticleDao) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.UTime = now
	// gorm忽略零值的特性, 会用主键更新
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id = ? AND author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  article.UTime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能作者非法, id: %d, author_id %d", article.Id, article.AuthorId)
	}
	return nil
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(255)"`
	Content string `gorm:"type=BLOB"`

	// 按照这个索引, 创建时间倒序排序
	// 最佳实践是 在 AuthorId 和 CTime 创建联合索引
	AuthorId int64 `gorm:"index=aid_ctime"`

	CTime int64 `gorm:"index=aid_ctime"`
	UTime int64
}
