package service

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/cronjob/domain"
)

type JobService interface {
	// Preempt 抢占
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
}
