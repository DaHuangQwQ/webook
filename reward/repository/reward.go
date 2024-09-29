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

func NewRewardRepository(dao dao.RewardDAO, cache cache.RewardCache) RewardRepository {
	return &rewardRepository{dao: dao, cache: cache}
}

func (repo *rewardRepository) CreateReward(ctx context.Context, reward domain.Reward) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(reward))
}

func (repo *rewardRepository) GetReward(ctx context.Context, rid int64) (domain.Reward, error) {
	res, err := repo.dao.GetReward(ctx, rid)
	return repo.toDomain(res), err
}

func (repo *rewardRepository) GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	return repo.cache.GetCachedCodeURL(ctx, r)
}

func (repo *rewardRepository) CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error {
	return repo.cache.CachedCodeURL(ctx, cu, r)
}

func (repo *rewardRepository) UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error {
	return repo.dao.UpdateStatus(ctx, rid, status.AsUint8())
}

func (repo *rewardRepository) toEntity(reward domain.Reward) dao.Reward {
	return dao.Reward{
		Id:        reward.Id,
		Biz:       reward.Target.Biz,
		BizId:     reward.Target.BizId,
		BizName:   reward.Target.BizName,
		TargetUid: reward.Target.Uid,
		Status:    reward.Status.AsUint8(),
		Uid:       reward.Uid,
		Amount:    reward.Amt,
	}
}

func (repo *rewardRepository) toDomain(reward dao.Reward) domain.Reward {
	return domain.Reward{
		Id: reward.Id,
		Target: domain.Target{
			Biz:     reward.Biz,
			BizId:   reward.BizId,
			BizName: reward.BizName,
			Uid:     reward.Uid,
		},
		Status: domain.RewardStatus(reward.Status),
		Uid:    reward.Uid,
		Amt:    reward.Amount,
	}
}
