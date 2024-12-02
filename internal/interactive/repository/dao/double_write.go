package dao

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
)

const (
	patternSrcOnly  = "SRC_ONLY"
	patternDstOnly  = "DST_ONLY"
	patternSrcFirst = "SRC_FIRST"
	patternDstFirst = "DST_FIRST"
)

type DoubleWriteDao struct {
	src     InteractiveDao
	dst     InteractiveDao
	pattern *atomicx.Value[string]
}

func NewDoubleWriteDao(src InteractiveDao, dst InteractiveDao) *DoubleWriteDao {
	return &DoubleWriteDao{src: src, dst: dst, pattern: atomicx.NewValueOf(patternSrcOnly)}
}

func (d *DoubleWriteDao) UpdatePattern(str string) {
	d.pattern.Store(str)
}

func (d *DoubleWriteDao) InteractiveReadCnt(ctx context.Context, biz string, bizId int64) error {
	switch d.pattern.Load() {
	case patternSrcOnly:
		return d.src.InteractiveReadCnt(ctx, biz, bizId)
	case patternSrcFirst:
		err := d.src.InteractiveReadCnt(ctx, biz, bizId)
		if err != nil {
			// src 都没写成功 等修复
			return err
		}
		// src 成功了 dst 失败了 等修复
		err = d.dst.InteractiveReadCnt(ctx, biz, bizId)
		if err != nil {
			// log
		}
		return nil
	case patternDstFirst:
		err := d.dst.InteractiveReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}
		err = d.src.InteractiveReadCnt(ctx, biz, bizId)
		if err != nil {
			// log
		}
		return nil
	case patternDstOnly:
		return d.dst.InteractiveReadCnt(ctx, biz, bizId)
	default:
		return errors.New("未知错误")
	}
}

func (d *DoubleWriteDao) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDao) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDao) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDao) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDao) GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDao) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	switch d.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return d.src.Get(ctx, biz, id)
	case patternDstOnly, patternDstFirst:
		return d.dst.Get(ctx, biz, id)
	default:
		return Interactive{}, errors.New("未知错误")
	}
}

func (d *DoubleWriteDao) GetByIds(ctx context.Context, biz string, bizIds []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}
