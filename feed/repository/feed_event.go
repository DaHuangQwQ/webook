package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"webook/feed/domain"
	"webook/feed/repository/cache"
	"webook/feed/repository/dao"
)

var FolloweesNotFound = cache.FolloweesNotFound

type FeedEventRepo interface {
	// CreatePushEvents 批量推事件
	CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error
	// CreatePullEvent 创建拉事件
	CreatePullEvent(ctx context.Context, event domain.FeedEvent) error
	// FindPullEvents 获取拉事件，也就是关注的人发件箱里面的事件
	FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error)
	// FindPushEvents 获取推事件，也就是自己收件箱里面的事件
	FindPushEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
	// FindPullEventsWithTyp 获取某个类型的拉事件，
	FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error)
	// FindPushEvents 获取某个类型的推事件，也就
	FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error)
}

type feedEventRepo struct {
	pullDao   dao.FeedPullEventDAO
	pushDao   dao.FeedPushEventDAO
	feedCache cache.FeedEventCache
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO, feedCache cache.FeedEventCache) FeedEventRepo {
	return &feedEventRepo{
		pullDao:   pullDao,
		pushDao:   pushDao,
		feedCache: feedCache,
	}
}

func (f *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventListWithTyp(ctx, typ, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.GetPushEventsWithTyp(ctx, typ, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) SetFollowees(ctx context.Context, follower int64, followees []int64) error {
	return f.feedCache.SetFollowees(ctx, follower, followees)
}

func (f *feedEventRepo) GetFollowees(ctx context.Context, follower int64) ([]int64, error) {
	followees, err := f.feedCache.GetFollowees(ctx, follower)
	if errors.Is(err, cache.FolloweesNotFound) {
		return nil, FolloweesNotFound
	}
	return followees, err
}

func (f *feedEventRepo) FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pullDao.FindPullEventList(ctx, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) FindPushEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := f.pushDao.GetPushEvents(ctx, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func (f *feedEventRepo) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	pushEvents := make([]dao.FeedPushEvent, 0, len(events))
	for _, e := range events {
		pushEvents = append(pushEvents, convertToPushEventDao(e))
	}
	return f.pushDao.CreatePushEvents(ctx, pushEvents)
}

func (f *feedEventRepo) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, convertToPullEventDao(event))
}

func convertToPushEventDao(event domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		CTime:   event.CTime.Unix(),
	}
}

func convertToPullEventDao(event domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		CTime:   event.CTime.Unix(),
	}

}

func convertToPushEventDomain(event dao.FeedPushEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		CTime: time.Unix(event.CTime, 0),
		Ext:   ext,
	}
}

func convertToPullEventDomain(event dao.FeedPullEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		CTime: time.Unix(event.CTime, 0),
		Ext:   ext,
	}
}
