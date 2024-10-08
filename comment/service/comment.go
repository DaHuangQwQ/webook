package service

import (
	"context"
	"webook/comment/domain"
	"webook/comment/repository"
)

type commentService struct {
	repo repository.CommentRepository
}

func (c *commentService) GetMoreReplies(ctx context.Context,
	rid int64,
	maxID int64, limit int64) ([]domain.Comment, error) {
	return c.repo.GetMoreReplies(ctx, rid, maxID, limit)
}

func NewCommentSvc(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

func (c *commentService) GetCommentList(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]domain.Comment, error) {
	list, err := c.repo.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	return list, err
}

func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	return c.repo.DeleteComment(ctx, domain.Comment{
		Id: id,
	})
}

func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.repo.CreateComment(ctx, comment)
}
