package dao

import (
	"database/sql"
	"gorm.io/gorm"
)

type Comment struct {
	Id int64 `gorm:"autoIncrement,primaryKey"`
	// 发表评论的人
	// 也就是说，如果你需要查询某个人发表的所有的评论，那么你需要在这里创建一个索引
	Uid int64
	// 被评价的东西
	// 这里要不要建索引？
	Biz     string `gorm:"index:biz_type_id"`
	BizID   int64  `gorm:"index:biz_type_id"`
	Content string

	// 我的根评论是哪个
	// 也就是说，如果这个字段是 NULL，它是根评论
	RootID sql.NullInt64 `gorm:"column:root_id;index"`

	// 这个是 NULL，也是根评论
	PID sql.NullInt64 `gorm:"column:pid;index"`
	// 级联删除
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`

	CTime int64
	// 事实上，大部分平台是不允许修改评论的
	UTime int64
}

func (*Comment) TableName() string {
	return "comments"
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&Comment{})
}
