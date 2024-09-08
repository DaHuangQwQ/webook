package dao

import (
	"context"
	"gorm.io/gorm"
)

type RecruitmentDao interface {
	DeleteByIds(ctx context.Context, ids []int64) error
	Create(ctx context.Context, recruitment Recruitment) error
}

type GormRecruitmentDao struct {
	db *gorm.DB
}

func NewGormRecruitmentDao(db *gorm.DB) RecruitmentDao {
	return &GormRecruitmentDao{db: db}
}

func (dao *GormRecruitmentDao) DeleteByIds(ctx context.Context, ids []int64) error {
	return dao.db.WithContext(ctx).Where("id in (?)", ids).Delete(&Recruitment{}).Error
}

func (dao *GormRecruitmentDao) Create(ctx context.Context, recruitment Recruitment) error {
	return dao.db.WithContext(ctx).Create(&recruitment).Error
}

type Recruitment struct {
	Id          uint64 `gorm:"primaryKey;autoIncrement;comment:'id'" json:"id"` // id (bigint)
	Name        string `gorm:"type:varchar(10);comment:'姓名'" json:"name"`       // 姓名 (varchar(10))
	StudentID   string `gorm:"type:varchar(10);comment:'学号'" json:"student_id"` // 学号 (varchar(10))
	Major       uint8  `gorm:"type:tinyint;comment:'专业'" json:"major"`          // 专业 (tinyint)
	Situation   string `gorm:"type:text;comment:'基础情况'" json:"situation"`       // 基础情况 (text)
	Expectation string `gorm:"type:text;comment:'未来期望'" json:"expectation"`     // 未来期望 (text)
	Selfie      string `gorm:"type:varchar(255);comment:'自拍图片'" json:"selfie"`  // 自拍图片 (varchar(255))
	ErrorNum    int    `gorm:"type:int;comment:'错误次数'" json:"error_num"`        // 错误次数 (int)
	UTime       int64  `gorm:"autoUpdateTime;comment:'更新时间'" json:"u_time"`     // 更新时间 (datetime)
	CTime       int64  `gorm:"autoCreateTime;comment:'创建时间'" json:"c_time"`     // 创建时间 (datetime)
}
