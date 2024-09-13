package job

import (
	"context"
	"time"
	"webook/internal/service"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
	}
}

func (r *RankingJob) Name() string {
	return "RankingJob"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}
