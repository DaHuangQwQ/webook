package repository

import (
	"context"
	"fmt"
	"github.com/DaHuangQwQ/gutil/slice"
	"github.com/spf13/viper"
	"time"
	"webook/article/domain"
	"webook/article/repository/cache"
	"webook/article/repository/dao"
	"webook/bff/api"
	"webook/internal/repository/dao/oss"
)

type CachedArticleRepository struct {
	dao dao.ArticleDao
	oss oss.Client

	cache cache.ArticleCache
}

func NewCachedArticleRepository(dao dao.ArticleDao, oss oss.Client, cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		oss:   oss,
		cache: cache,
	}
}

func (c *CachedArticleRepository) ListPub(ctx context.Context, start time.Time, pageNum int, pageSize int) ([]domain.Article, error) {
	arts, err := c.dao.ListPub(ctx, start, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.PublishedArticle, domain.Article](arts,
		func(idx int, src dao.PublishedArticle) domain.Article {
			return c.toPubDomain(src)
		}), nil
}

func (c *CachedArticleRepository) DeleteByIds(ctx context.Context, ids []int64) error {
	return c.dao.DeleteByIds(ctx, ids)
}

func (c *CachedArticleRepository) ListAll(ctx context.Context, req api.PageReq) ([]domain.Article, error) {
	res, err := c.dao.ListAll(ctx, req.PageNum, req.PageSize)
	articles := make([]domain.Article, len(res))
	for i, v := range res {
		articles[i] = c.toPubDomain(v)
		articles[i].Content = articles[i].Abstract()
	}
	return articles, err
}

func (c *CachedArticleRepository) GetById(ctx context.Context, articleId int64) (domain.Article, error) {
	res, err := c.dao.FindById(ctx, articleId)
	return c.toPubDomain(res), err
}

func (c *CachedArticleRepository) ListById(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				c.preCache(ctx, data)
			}()
			return data[:limit], err
		}
	}
	res, err := c.dao.GetListByAuthor(ctx, uid, offset, limit)
	articles := make([]domain.Article, len(res))
	for i, v := range res {
		articles[i] = c.toDomain(v)
	}
	go func() {
		err := c.cache.SetFirstPage(ctx, uid, articles)
		if err != nil {
			// 日志
			return
		}
		c.preCache(ctx, articles)
	}()
	return articles, err
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	defer func() {
		err := c.cache.DelFirstPage(ctx, article.Id)
		if err != nil {
			// log
			return
		}
	}()
	return c.dao.Insert(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Img_Update(ctx context.Context, file []byte, fileType string) (string, error) {
	type Config struct {
		ENDP string `yaml:"ENDP"`
	}
	var config Config
	err := viper.UnmarshalKey("OSS", &config)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("articleImg/%d.%s", time.Now().UnixMilli(), fileType)
	ossAdress := "https://ceit." + config.ENDP + "/" + fileName
	return ossAdress, c.oss.UploadFile(ctx, fileName, file)
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	defer func() {
		err := c.cache.DelFirstPage(ctx, article.Id)
		if err != nil {
			// log
			return
		}
	}()
	return c.dao.UpdateById(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(article))
	if err == nil {
		er := c.cache.DelFirstPage(ctx, article.Author.Id)
		if er != nil {
			// 也要记录日志
		}
	}
	// 在这里尝试，设置缓存
	go func() {
		ct, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		article.Author = domain.Author{
			Id: article.Author.Id,
		}
		er := c.cache.SetPub(ct, article)
		if er != nil {
			// 记录日志
		}
	}()
	return id, err
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, articleId int64, authorId int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, articleId, authorId, status.ToUInt8())
}

func (c *CachedArticleRepository) GetList(ctx context.Context) ([]domain.Article, error) {
	articlesList, err := c.dao.GetList(ctx)
	if err != nil {
		return nil, err
	}
	articlesDomainList := make([]domain.Article, len(articlesList))
	for index, article := range articlesList {
		articlesDomainList[index] = domain.Article{
			Id:      article.Id,
			Title:   article.Title,
			Content: article.Content,
			ImgUrl:  article.ImgUrl,
			Type:    article.Type,
			Author: domain.Author{
				Id: article.AuthorId,
			},
		}
	}
	return articlesDomainList, nil
}

func (c *CachedArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUInt8(),
		ImgUrl:   article.ImgUrl,
		Type:     article.Type,
	}
}

func (c *CachedArticleRepository) toDomain(article dao.Article) domain.Article {
	return domain.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			Id: article.AuthorId,
		},
		Status: domain.ArticleStatus(article.Status),
		ImgUrl: article.ImgUrl,
		Type:   article.Type,
	}
}

func (c *CachedArticleRepository) toPubDomain(article dao.PublishedArticle) domain.Article {
	return domain.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author: domain.Author{
			Id: article.AuthorId,
		},
		Status: domain.ArticleStatus(article.Status),
		ImgUrl: article.ImgUrl,
		Type:   article.Type,
	}
}

func (c *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const size = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < size {
		err := c.cache.Set(ctx, arts[0])
		if err != nil {
			// 记录缓存
		}
	}
}
