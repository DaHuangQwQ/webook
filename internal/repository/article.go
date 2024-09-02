package repository

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"time"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/repository/dao/oss"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, articleId int64, authorId int64, status domain.ArticleStatus) error
	GetList(ctx context.Context) ([]domain.Article, error)
	ListById(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, uid int64, articleId int64) (domain.Article, error)

	Img_Update(ctx context.Context, file []byte, fileType string) (string, error)
	ListAll(ctx context.Context, req api.PageReq) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao
	oss oss.Client

	userRepo UserRepository

	cache cache.ArticleCache
}

func NewCachedArticleRepository(dao dao.ArticleDao, oss oss.Client, userRepo UserRepository, cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		oss:      oss,
		userRepo: userRepo,
		cache:    cache,
	}
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

func (c *CachedArticleRepository) GetById(ctx context.Context, uid, articleId int64) (domain.Article, error) {
	res, err := c.dao.FindById(ctx, uid, articleId)
	return c.toDomain(res), err
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		user, er := c.userRepo.FindByID(ctx, article.Author.Id)
		if er != nil {
			// 要记录日志
			return
		}
		article.Author = domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		}
		er = c.cache.SetPub(ctx, article)
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
		author, err := c.userRepo.FindByID(ctx, article.AuthorId)
		if err != nil {
			return nil, err
		}
		articlesDomainList[index] = domain.Article{
			Id:      article.Id,
			Title:   article.Title,
			Content: article.Content,
			ImgUrl:  article.ImgUrl,
			Type:    article.Type,
			Author: domain.Author{
				Id:   author.Id,
				Name: author.Nickname,
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
