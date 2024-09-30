package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type GormJobDAO struct {
	db *gorm.DB
}

func NewGormJobDAO(db *gorm.DB) JobDao {
	return &GormJobDAO{
		db: db,
	}
}

func (dao *GormJobDAO) Release(ctx context.Context, id int64) error {
	// TODO where version = ?
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"u_time": time.Now().UnixMilli(),
	}).Error
}

func (dao *GormJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		now := time.Now().UnixMilli()
		//分布式任务调度系统
		//1. 一次取一批， 一次性取出100个，然后随机从某一条开始抢占
		//2. 随机偏移量，兜底：第一次没查到，偏移量回归到 0
		//3. id取余分配， 兜底不加余数条件

		var job Job
		err := db.Where("status = ? AND next_time <= ?", jobStatusRunning, now).First(&job).Error
		if err != nil {
			return Job{}, err
		}
		// 乐观锁 CAS操作，compare and swap
		// 用乐观锁 取代 for update
		res := db.Model(&Job{}).Where("id = ? AND version = ?", job.Id, job.Version).Updates(map[string]any{
			"status":  jobStatusRunning,
			"version": job.Version + 1,
			"u_time":  now,
			"c_time":  now,
		})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			// 抢占失败
			continue
		}
		return job, err
	}
}

func (dao *GormJobDAO) UpdateUTime(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"u_time": time.Now().UnixMilli(),
	}).Error
}

func (dao *GormJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (dao *GormJobDAO) Stop(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"u_time": time.Now().UnixMilli(),
	}).Error
}

type Job struct {
	Id       int64 `gorm:"primary_key,AUTO_INCREMENT"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`
	Status   int
	Version  int64
	Cron     string
	NextTime int64 `gorm:"index"`
	UTime    int64
	CTime    int64
}
