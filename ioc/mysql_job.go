package ioc

import (
	"context"
	"github.com/DaHuangQwQ/gpkg/logger"
	"github.com/DaHuangQwQ/webook/internal/cronjob/domain"
	"github.com/DaHuangQwQ/webook/internal/cronjob/service"
	service2 "github.com/DaHuangQwQ/webook/internal/ranking/service"
	"github.com/DaHuangQwQ/webook/internal_temp/job"
	"time"
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
