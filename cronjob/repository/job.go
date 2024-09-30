package repository

import (
	"context"
	"time"
	"webook/cronjob/domain"
	"webook/cronjob/repository/dao"
)

type PreemptCronJobRepository struct {
	dao dao.JobDao
}

func NewPreemptCronJobRepository(dao dao.JobDao) JobRepository {
	return &PreemptCronJobRepository{
		dao: dao,
	}
}

func (repo *PreemptCronJobRepository) UpdateUTime(ctx context.Context, id int64) error {
	return repo.dao.UpdateUTime(ctx, id)
}

func (repo *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return repo.dao.UpdateNextTime(ctx, id, next)
}

func (repo *PreemptCronJobRepository) Stop(ctx context.Context, id int64) error {
	return repo.dao.Stop(ctx, id)
}

func (repo *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	return repo.dao.Release(ctx, id)
}

func (repo *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := repo.dao.Preempt(ctx)
	return repo.toDomain(job), err
}

func (repo *PreemptCronJobRepository) toDomain(job dao.Job) domain.Job {
	return domain.Job{
		Id:       job.Id,
		Name:     job.Name,
		Cfg:      job.Cfg,
		Executor: job.Executor,
	}
}
