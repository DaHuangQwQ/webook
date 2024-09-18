package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
	"webook/internal/service"
	"webook/pkg/logger"
)

type RankingJob struct {
	svc     service.RankingService
	l       logger.LoggerV1
	timeout time.Duration
	client  *rlock.Client
	key     string

	localLock *sync.Mutex
	lock      *rlock.Lock

	// 作业提示
	// 随机生成一个，就代表当前负载。你可以每隔一分钟生成一个
	load int32
}

func NewRankingJob(
	svc service.RankingService,
	l logger.LoggerV1,
	client *rlock.Client,
	timeout time.Duration) *RankingJob {
	return &RankingJob{svc: svc,
		key:       "job:ranking",
		l:         l,
		client:    client,
		localLock: &sync.Mutex{},
		timeout:   timeout}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// go fun() { r.Run()}

// 释放分布式锁失败 怎么办 重试，及时告警 人手工处理
// 能控制同一时刻 只有一个 goroutine 执行， 无法保证 运行完成后 其他实例接着执行

func (r *RankingJob) Run() error {
	r.localLock.Lock()
	lock := r.lock
	if lock == nil {
		// 抢分布式锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout,
			&rlock.FixIntervalRetry{
				Interval: time.Millisecond * 100,
				Max:      3,
				// 重试的超时
			}, time.Second)
		if err != nil {
			r.localLock.Unlock()
			r.l.Warn("获取分布式锁失败", logger.Error(err))
			return nil
		}
		r.lock = lock
		r.localLock.Unlock()
		go func() {
			// 并不是非得一半就续约
			er := lock.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				// 续约失败了
				// 你也没办法中断当下正在调度的热榜计算（如果有）
			}
			r.localLock.Lock()
			r.lock = nil
			//lock.Unlock()
			r.localLock.Unlock()
		}()
	}
	// 这边就是你拿到了锁
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.lock = nil
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}
