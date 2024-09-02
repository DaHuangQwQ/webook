package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	InteractiveReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
}

type GormInteractiveDao struct {
	db *gorm.DB
}

// Get
// 拿出阅读数, 点赞数, 收藏数
func (dao *GormInteractiveDao) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&res).Error
	return res, err
}

func NewGormInteractiveDao(db *gorm.DB) InteractiveDao {
	return &GormInteractiveDao{
		db: db,
	}
}

func (dao *GormInteractiveDao) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&res).Error
	return res, err
}

func (dao *GormInteractiveDao) GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&res).Error
	return res, err
}

func (dao *GormInteractiveDao) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.CTime = now
	cb.UTime = now
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"u_time":      now,
			}),
		}).Create(&Interactive{
			Biz:        cb.Biz,
			BizId:      cb.BizId,
			CollectCnt: 1,
			CTime:      now,
			UTime:      now,
		}).Error
	})
}

func (dao *GormInteractiveDao) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"u_time": now,
				"status": 1,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Biz:    biz,
			BizId:  id,
			Status: 1,
			UTime:  now,
			CTime:  now,
		}).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"u_time":   now,
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   id,
			LikeCnt: 1,
			CTime:   now,
			UTime:   now,
		}).Error
	})
}

func (dao *GormInteractiveDao) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid=? AND biz_id = ? AND biz=?", uid, id, biz).
			Updates(map[string]interface{}{
				"u_time": now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).
			Where("biz =? AND biz_id=?", biz, id).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"u_time":   now,
			}).Error
	})
}

// InteractiveReadCnt new or update
func (dao *GormInteractiveDao) InteractiveReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"read_cnt": gorm.Expr("read_cnt + 1"),
				"u_time":   time.Now().UnixMilli(),
			}),
		}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		UTime:   now,
		CTime:   now,
	}).Error
}

type Interactive struct {
	Id int64 `json:"id" gorm:"primary_key;autoIncrement"`

	// 联合唯一索引
	BizId int64  `json:"biz_id" gorm:"uniqueIndex:uk_biz_id_biz"`
	Biz   string `json:"biz" gorm:"type:varchar(128);uniqueIndex:uk_biz_id_biz"`

	// 阅读计数
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
	CTime      int64 `json:"c_time"`
	UTime      int64 `json:"u_time"`
}

type UserLikeBiz struct {
	Id    int64  `json:"id" gorm:"primary_key;autoIncrement"`
	Uid   int64  `json:"uid"    gorm:"uniqueIndex:uk_uid_biz_id_biz"`
	Biz   string `json:"biz"    gorm:"type:varchar(128);uniqueIndex:uk_uid_biz_id_biz"`
	BizId int64  `json:"biz_id" gorm:"uniqueIndex:uk_uid_biz_id_biz"`
	// 软删除
	Status uint8 `json:"status" gorm:"type:tinyint(1)"`

	CTime int64 `json:"c_time"`
	UTime int64 `json:"u_time"`
}

type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 这边还是保留了了唯一索引
	Uid   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	// 收藏夹的ID
	// 收藏夹ID本身有索引
	Cid   int64 `gorm:"index"`
	UTime int64
	CTime int64
}
