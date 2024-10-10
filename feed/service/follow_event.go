package service

import (
	"context"
	"time"
	"webook/feed/domain"
	"webook/feed/repository"
)

const (
	FollowEventName = "follow_event"
)

type FollowEventHandler struct {
	repo repository.FeedEventRepo
}

func NewFollowEventHandler(repo repository.FeedEventRepo) Handler {
	return &FollowEventHandler{
		repo: repo,
	}
}

func (f *FollowEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	return f.repo.FindPushEventsWithTyp(ctx, FollowEventName, uid, timestamp, limit)
}

// CreateFeedEvent 创建跟随方式
// 如果 A 关注了 B，那么
// follower 就是 A
// followee 就是 B
func (f *FollowEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	followee, err := ext.Get("followee").AsInt64()
	if err != nil {
		return err
	}
	return f.repo.CreatePushEvents(ctx, []domain.FeedEvent{{
		Uid:   followee,
		Type:  FollowEventName,
		CTime: time.Now(),
		Ext:   ext,
	}})
}
