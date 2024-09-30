package service

import (
	"context"
	"time"
	"webook/cronjob/domain"
	"webook/cronjob/repository"
	"webook/pkg/logger"
)

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.LoggerV1
}

func NewCronJobService(repo repository.JobRepository, l logger.LoggerV1) JobService {
	return &CronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second,
	}
}

func (svc *CronJobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	next := job.NextTime()
	if next.IsZero() {
		// 没有下一次
		return svc.repo.Stop(ctx, job.Id)
	}
	return svc.repo.UpdateNextTime(ctx, job.Id, next)
}

func (svc *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := svc.repo.Preempt(ctx)

	ticker := time.NewTicker(svc.refreshInterval)
	go func() {
		for range ticker.C {
			svc.refresh(job.Id)
		}
	}()
	job.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return svc.repo.Release(ctx, job.Id)
	}
	return job, err
}

// refresh 续约
func (svc *CronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := svc.repo.UpdateUTime(ctx, id)
	if err != nil {
		// 考虑重试
		svc.l.Error("续约失败", logger.Error(err), logger.Int64("jid", id))
	}
}
