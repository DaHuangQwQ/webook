package job

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
	"webook/pkg/logger"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, job domain.Job) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, job domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{funcs: make(map[string]func(ctx context.Context, job domain.Job) error)}
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

func (l *LocalFuncExecutor) Name() string {
	return "LocalFuncExecutor"
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.Job) error {
	fn, ok := l.funcs[job.Name]
	if !ok {
		return fmt.Errorf("未知任务%s 你是否注册了", job.Name)
	}
	return fn(ctx, job)
}

// Scheduler 调度器
type Scheduler struct {
	execs   map[string]Executor
	svc     service.JobService
	l       logger.LoggerV1
	limiter *semaphore.Weighted
}

func NewScheduler(svc service.JobService, l logger.LoggerV1) *Scheduler {
	return &Scheduler{
		svc:     svc,
		l:       l,
		limiter: semaphore.NewWeighted(100),
		execs:   make(map[string]Executor),
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := s.limiter.Acquire(ctx, 1)

		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		job, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			s.l.Error("抢占任务失败", logger.Error(err))
		}

		exec, ok := s.execs[job.Executor]
		if !ok {
			s.l.Error("未找到对应的执行器")
			continue
		}

		go func() {
			s.limiter.Release(1)
			// 异步调度 不阻塞主调度循环
			defer func() {
				er := job.CancelFunc()
				if er != nil {
					s.l.Error("释放任务失败", logger.Error(err))
				}
			}()
			er := exec.Exec(ctx, job)
			if er != nil {
				s.l.Error("任务执行失败", logger.Error(er))
			}

			er = s.svc.ResetNextTime(ctx, job)
			if er != nil {
				s.l.Error("设置下一次执行时间失败", logger.Error(er))
			}
		}()

	}
}
