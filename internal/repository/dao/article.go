package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Upsert(ctx context.Context, article PublishedArticle) error
	Sync(ctx context.Context, article Article) (int64, error)
	SyncStatus(ctx context.Context, articleID int64, authorId int64, status uint8) error
	GetList(ctx context.Context) ([]PublishedArticle, error)
	FindById(ctx context.Context, uid int64, articleId int64) (Article, error)
	GetListByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	ListAll(ctx context.Context, PageNum int, PageSize int) ([]PublishedArticle, error)
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewGormArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

func (dao *GormArticleDao) ListAll(ctx context.Context, PageNum int, PageSize int) ([]PublishedArticle, error) {
	var articles []PublishedArticle
	err := dao.db.WithContext(ctx).Offset(PageNum).Limit(PageSize).Order("id desc").Find(&articles).Error
	return articles, err
}

func (dao *GormArticleDao) GetListByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var articles []Article
	err := dao.db.WithContext(ctx).Where("author_id = ?", uid).Offset(offset).Limit(limit).Find(&articles).Error
	return articles, err
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
			"status":  article.Status,
			"img_url": article.ImgUrl,
			"type":    article.Type,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能作者非法, id: %d, author_id %d", article.Id, article.AuthorId)
	}
	return nil
}

func (dao *GormArticleDao) FindById(ctx context.Context, uid int64, articleId int64) (Article, error) {
	article := Article{}
	err := dao.db.WithContext(ctx).Model(&Article{}).Where("id = ? AND uid = ?", articleId, uid).
		First(&article).Error
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

func (dao *GormArticleDao) Sync(ctx context.Context, article Article) (int64, error) {
	// 先操作 制作库 后操作 线上库
	var (
		id = article.Id
	)
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
			//articleInfo Article
		)
		txDao := NewGormArticleDao(tx)
		if id > 0 {
			//articleInfo, err = txDao.FindById(ctx, article.Id, article.AuthorId)
			err = txDao.UpdateById(ctx, article)
			if err != nil {
				return err
			}
		} else {
			id, err = txDao.Insert(ctx, article)
			if err != nil {
				return err
			}
		}
		err = txDao.Upsert(ctx, PublishedArticle{Article: article})
		return err
	})
	return id, err
}

func (dao *GormArticleDao) SyncStatus(ctx context.Context, articleID int64, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", articleID, authorId).
			Updates(map[string]any{
				"status": status,
				"u_time": now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("作者不对, 有人在搞, article_id: %d, author_id %d", articleID, authorId)
		}

		return tx.Model(&PublishedArticle{}).Where("id = ?", articleID).
			Updates(map[string]any{
				"status": status,
				"u_time": now,
			}).Error

	})
}

func (dao *GormArticleDao) Upsert(ctx context.Context, article PublishedArticle) error {
	// insert or update
	now := time.Now().UnixMilli()
	article.UTime = now
	article.CTime = now
	// OnConflict 数据冲突了
	err := dao.db.Clauses(clause.OnConflict{
		// mysql 只关心这些 不遵循 sql2003 标准
		DoUpdates: clause.Assignments(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"u_time":  now,
			"status":  article.Status,
			"img_url": article.ImgUrl,
			"type":    article.Type,
		}),
	}).Create(&article).Error
	return err
}

func (dao *GormArticleDao) GetList(ctx context.Context) ([]PublishedArticle, error) {
	res := make([]PublishedArticle, 0)
	err := dao.db.WithContext(ctx).Find(&res).Error
	return res, err
}

type PublishedArticle struct {
	Article
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(255)"`
	Content string `gorm:"type=BLOB"`

	// 按照这个索引, 创建时间倒序排序
	// 最佳实践是 在 AuthorId 和 CTime 创建联合索引
	AuthorId int64  `gorm:"index=aid_ctime"`
	ImgUrl   string `json:"img_url" gorm:"varchar(255)"`
	Type     string `json:"type" gorm:"varchar(255)"`

	Status uint8

	CTime int64 `gorm:"index=aid_ctime"`
	UTime int64
}
