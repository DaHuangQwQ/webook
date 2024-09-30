package service

import (
	"context"
	"webook/cronjob/domain"
)

type JobService interface {
	// Preempt 抢占
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
}
