package dao

import (
	"context"
	"gorm.io/gorm"
)

type GORMCommentDAO struct {
	db *gorm.DB
}

func (c *GORMCommentDAO) FindRepliesByRid(ctx context.Context,
	rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}

func (c *GORMCommentDAO) FindOneByIDs(ctx context.Context, ids []int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("id in ?", ids).
		First(&res).
		Error
	return res, err
}

func (c *GORMCommentDAO) FindByBiz(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

// FindRepliesByPid 查找评论的直接评论
func (c *GORMCommentDAO) FindRepliesByPid(ctx context.Context,
	pid int64,
	offset,
	limit int) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (c *GORMCommentDAO) Insert(ctx context.Context, u Comment) error {
	return c.db.
		WithContext(ctx).
		Create(u).
		Error
}

func (c *GORMCommentDAO) FindCommentList(ctx context.Context, u Comment) ([]Comment, error) {
	var res []Comment
	builder := c.db.WithContext(ctx)
	if u.Id == 0 {
		builder = builder.
			Where("biz=?", u.Biz).
			Where("biz_id=?", u.BizID).
			Where("root_id is null")
	} else {
		builder = builder.Where("root_id=? or id =?", u.Id, u.Id)
	}
	err := builder.Find(&res).Error
	return res, err

}

func (c *GORMCommentDAO) Delete(ctx context.Context, u Comment) error {
	return c.db.WithContext(ctx).Delete(&Comment{
		Id: u.Id,
	}).Error
}
