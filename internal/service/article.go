package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (id int64, err error)
	Img_Update(ctx context.Context, file multipart.File, fileType string) (string, error)
	Withdraw(ctx context.Context, article domain.Article) error
	Publish(ctx context.Context, article domain.Article) (int64, error)
	GetList(ctx context.Context) (list []domain.Article, err error)
	List(ctx context.Context, req api.PageReq) (list []domain.Article, err error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (s *articleService) List(ctx context.Context, req api.PageReq) (list []domain.Article, err error) {
	return s.repo.ListAll(ctx, req)
}

func (s *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := s.repo.Update(ctx, article)
		return article.Id, err
	}
	return s.repo.Create(ctx, article)
}

func (s *articleService) Img_Update(ctx context.Context, file multipart.File, fileType string) (string, error) {
	fileBytes, err := s.fileToBytes(file)
	if err != nil {
		return " ", fmt.Errorf("%w", err)
	}
	return s.repo.Img_Update(ctx, fileBytes, fileType)
}

func (s *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// 制作库
	//var (
	//	id  = article.Id
	//	err error
	//)
	//if article.Id > 0 {
	//	err = s.repo.Update(ctx, article)
	//} else {
	//	id, err = s.repo.Create(ctx, article)
	//}
	//
	//article.Id = id
	// 线上库 俩库id应该是相等的

	// save: update or create
	article.Status = domain.ArticleStatusPublished
	return s.repo.Sync(ctx, article)
}

func (s *articleService) Withdraw(ctx context.Context, article domain.Article) error {
	return s.repo.SyncStatus(ctx, article.Id, article.Author.Id, domain.ArticleStatusPrivate)
}

func (s *articleService) GetList(ctx context.Context) (list []domain.Article, err error) {
	return s.repo.GetList(ctx)
}

func (s *articleService) fileToBytes(file multipart.File) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
