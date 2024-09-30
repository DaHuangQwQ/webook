package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"
	"webook/article/domain"
	"webook/article/events"
	"webook/article/repository"
	"webook/bff/api"
)

type articleService struct {
	repo     repository.ArticleRepository
	userRepo repository.AuthorRepository
	producer events.Producer
}

func NewArticleService(repo repository.ArticleRepository, userRepo repository.AuthorRepository, producer events.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
		userRepo: userRepo,
	}
}

func (s *articleService) ListPub(ctx context.Context, start time.Time, pageNum, pageSize int) ([]domain.Article, error) {
	return s.repo.ListPub(ctx, start, pageNum, pageSize)
}

func (s *articleService) DeleteByIds(ctx context.Context, ids []int64) error {
	return s.repo.DeleteByIds(ctx, ids)
}

func (s *articleService) GetPublishedById(ctx context.Context, uid, articleId int64) (domain.Article, error) {
	article, err := s.repo.GetById(ctx, articleId)
	if err == nil {
		go func() {
			er := s.producer.ProduceReadEvent(ctx, events.ReadEvent{
				Uid: uid,
				Aid: articleId,
			})
			if er != nil {
				// log
				return
			}
		}()
	}
	return article, err
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
