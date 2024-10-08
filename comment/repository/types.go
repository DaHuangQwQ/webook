package repository

import (
	"context"
	"webook/comment/domain"
)

type CommentRepository interface {
	// FindByBiz 根据 ID 倒序查找
	// 并且会返回每个评论的三条直接回复
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]domain.Comment, error)
	// DeleteComment 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, comment domain.Comment) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	// GetCommentByIds 获取单条评论 支持批量获取
	GetCommentByIds(ctx context.Context, id []int64) ([]domain.Comment, error)
	GetMoreReplies(ctx context.Context, rid int64, id int64, limit int64) ([]domain.Comment, error)
}
