package article

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
	GetList(ctx context.Context) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao

	userRepo repository.UserRepository
}

func NewCachedArticleRepository(dao dao.ArticleDao, userRepo repository.UserRepository) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		userRepo: userRepo,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, c.toEntity(article))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(article))
	//if err == nil {
	//	er := c.cache.DelFirstPage(ctx, article.Author.Id)
	//	if er != nil {
	//		// 也要记录日志
	//	}
	//}
	//// 在这里尝试，设置缓存
	//go func() {
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//	defer cancel()
	//
	//	user, er := c.userRepo.FindByID(ctx, article.Author.Id)
	//	if er != nil {
	//		// 要记录日志
	//		return
	//	}
	//	article.Author = domain.Author{
	//		Id:   user.Id,
	//		Name: user.Nickname,
	//	}
	//	er = c.cache.SetPub(ctx, article)
	//	if er != nil {
	//		// 记录日志
	//	}
	//}()
	return id, err
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
	}
}
