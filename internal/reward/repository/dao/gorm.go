package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type RewardGORMDAO struct {
	db *gorm.DB
}

func NewRewardGORMDAO(db *gorm.DB) RewardDAO {
	return &RewardGORMDAO{db: db}
}

func (dao *RewardGORMDAO) Insert(ctx context.Context, r Reward) (int64, error) {
	now := time.Now().UnixMilli()
	r.UTime = now
	r.CTime = now
	err := dao.db.WithContext(ctx).Create(&r).Error
	return r.Id, err
}

func (dao *RewardGORMDAO) GetReward(ctx context.Context, rid int64) (Reward, error) {
	// 通过 uid 来判定是自己的打赏，防止黑客捞数据
	var reward Reward
	err := dao.db.WithContext(ctx).Where("id = ?", rid).First(&reward).Error
	return reward, err
}

func (dao *RewardGORMDAO) UpdateStatus(ctx context.Context, rid int64, status uint8) error {
	return dao.db.WithContext(ctx).Model(&Reward{}).Where("id = ?", rid).Updates(map[string]any{
		"status": status,
		"u_time": time.Now().UnixMilli(),
	}).Error
}
