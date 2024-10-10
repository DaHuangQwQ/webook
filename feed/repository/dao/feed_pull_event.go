package dao

import (
	"context"
	"gorm.io/gorm"
)

// FeedPullEventDAO 拉模型
type FeedPullEventDAO interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
	FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

type FeedPullEvent struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	UID     int64  `gorm:"column:uid;type:int(11);not null;"`
	Type    string `gorm:"column:type;type:varchar(255);comment:类型"`
	Content string `gorm:"column:content;type:text;"`
	// 发生时间
	CTime int64 `gorm:"column:ctime;comment:发生时间"`
}

type feedPullEventDAO struct {
	db *gorm.DB
}

func NewFeedPullEventDAO(db *gorm.DB) FeedPullEventDAO {
	return &feedPullEventDAO{
		db: db,
	}
}

func (f *feedPullEventDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}

func (f *feedPullEventDAO) FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}
