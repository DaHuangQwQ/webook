package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	// Get
	// 从缓存里拿出阅读数, 点赞数, 收藏数
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
}

func NewCachedInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (repo *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := repo.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

// Get
// 从缓存里拿出阅读数, 点赞数, 收藏数
func (repo *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	res, err := repo.cache.Get(ctx, biz, id)
	if err == nil {
		return res, nil
	}
	// 查询数据库
	daoRes, err := repo.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	res = repo.toDomain(daoRes)
	go func() {
		er := repo.cache.Set(ctx, biz, id, res)
		if er != nil {
			// log
			return
		}
	}()
	return res, nil
}

func (repo *CachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	err := repo.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})
	go func() {
		err := repo.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
		if err != nil {
			// log
			return
		}
	}()
	return err
}

func (repo *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	// 插入点赞， 更新点赞计数， 更新缓存
	err := repo.dao.InsertLikeInfo(ctx, biz, id, uid)
	go func() {
		err := repo.cache.IncrLikeCntIfPresent(ctx, biz, id)
		if err != nil {
			// log
			return
		}
	}()
	return err
}

func (repo *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, id, uid)
	go func() {
		err := repo.cache.DecrLikeCntIfPresent(ctx, biz, id)
		if err != nil {
			// log
			return
		}
	}()
	return err
}

func (repo *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := repo.dao.InteractiveReadCnt(ctx, biz, bizId)
	go func() {
		err := repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
		if err != nil {
			// 记录一下
			return
		}
	}()
	return err
}

func (repo *CachedInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}
}
