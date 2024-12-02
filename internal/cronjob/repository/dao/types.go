package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

type JobDao interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUTime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&Job{})
}
