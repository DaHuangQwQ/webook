package repository

import (
	"context"
	"webook/reward/domain"
	"webook/reward/repository/cache"
	"webook/reward/repository/dao"
)

type rewardRepository struct {
	dao   dao.RewardDAO
	cache cache.RewardCache
}

func newRewardRepository(dao dao.RewardDAO, cache cache.RewardCache) RewardRepository {
	return &rewardRepository{dao: dao, cache: cache}
}

func (repo *rewardRepository) CreateReward(ctx context.Context, reward domain.Reward) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *rewardRepository) GetReward(ctx context.Context, rid int64) (domain.Reward, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *rewardRepository) GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *rewardRepository) CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error {
	//TODO implement me
	panic("implement me")
}

func (repo *rewardRepository) UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error {
	//TODO implement me
	panic("implement me")
}
