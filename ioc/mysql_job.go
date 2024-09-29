package ioc

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/job"
	"webook/internal/service"
	"webook/pkg/logger"
	service2 "webook/ranking/service"
)

func InitScheduler(
	l logger.LoggerV1,
	local job.Executor,
	svc service.JobService,
) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExecutor(svc service2.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	// 要在数据库 插入一条记录 ranking job
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		return svc.TopN(ctx)
	})
	return res
}
