package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/article"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (id int64, err error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo article.ArticleRepository
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (s *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := s.repo.Update(ctx, article)
		return article.Id, err
	}
	return s.repo.Create(ctx, article)
}

func (s *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// 制作库
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		err = s.repo.Update(ctx, article)
	} else {
		id, err = s.repo.Create(ctx, article)
	}

	article.Id = id
	// 线上库 俩库id应该是相等的

	// save: update or create
	return id, err
}