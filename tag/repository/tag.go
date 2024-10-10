package repository

import (
	"context"
	"github.com/DaHuangQwQ/gutil/slice"
	"time"
	"webook/pkg/logger"
	"webook/tag/domain"
	"webook/tag/repository/cache"
	"webook/tag/repository/dao"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag domain.Tag) (int64, error)
	BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
	GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error)
	GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error)
}

type CachedTagRepository struct {
	dao   dao.TagDAO
	cache cache.TagCache
	l     logger.LoggerV1
}

// PreloadUserTags 在 toB 的场景下，你可以提前预加载缓存, 全量缓存
func (repo *CachedTagRepository) PreloadUserTags(ctx context.Context) error {
	offset := 0
	const batch = 100
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		tags, err := repo.dao.GetTags(dbCtx, offset, batch)
		cancel()
		if err != nil {
			return err
		}
		// 你可以开 goroutine 来并发
		for _, tag := range tags {
			cCtx, cancel := context.WithTimeout(ctx, time.Second)
			err = repo.cache.Append(cCtx, tag.Uid, repo.toDomain(tag))
			cancel()
			if err != nil {
				return err
			}
		}
		if len(tags) < batch {
			break
		}
		offset = offset + batch
	}
	return nil
}

func (repo *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	tags, err := repo.dao.GetTagsById(ctx, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	return repo.dao.CreateTagBiz(ctx, slice.Map(tags, func(idx int, src int64) dao.TagBiz {
		return dao.TagBiz{
			Tid:   src,
			BizId: bizId,
			Biz:   biz,
			Uid:   uid,
		}
	}))
}

func (repo *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := repo.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}
	tags, err := repo.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	res = slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	})
	err = repo.cache.Append(ctx, uid, res...)
	if err != nil {
		// 记录日志
	}
	return res, nil
}

func (repo *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	tags, err := repo.dao.GetTagsByBiz(ctx, uid, biz, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	id, err := repo.dao.CreateTag(ctx, repo.toEntity(tag))
	if err != nil {
		return 0, err
	}
	// 也可以考虑用 DelTags
	err = repo.cache.Append(ctx, tag.Uid, tag)
	if err != nil {
		// 记录日志
	}
	return id, nil
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.LoggerV1) TagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}

func (repo *CachedTagRepository) toDomain(tag dao.Tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func (repo *CachedTagRepository) toEntity(tag domain.Tag) dao.Tag {
	return dao.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}
