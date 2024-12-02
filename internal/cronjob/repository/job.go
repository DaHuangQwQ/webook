package repository

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/cronjob/domain"
	dao2 "github.com/DaHuangQwQ/webook/internal/cronjob/repository/dao"
	"time"
)

type PreemptCronJobRepository struct {
	dao dao2.JobDao
}

func NewPreemptCronJobRepository(dao dao2.JobDao) JobRepository {
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

func (repo *PreemptCronJobRepository) toDomain(job dao2.Job) domain.Job {
	return domain.Job{
		Id:       job.Id,
		Name:     job.Name,
		Cfg:      job.Cfg,
		Executor: job.Executor,
	}
}
