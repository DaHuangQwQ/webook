package service

import (
	"context"
	"webook/reward/domain"
	"webook/reward/repository"
)

type WechatNativeRewardService struct {
	repo repository.RewardRepository
}

func NewWechatNativeRewardService(repo repository.RewardRepository) RewardService {
	return &WechatNativeRewardService{repo: repo}
}

func (w WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	//TODO implement me
	panic("implement me")
}

func (w WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	//TODO implement me
	panic("implement me")
}

func (w WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	//TODO implement me
	panic("implement me")
}
