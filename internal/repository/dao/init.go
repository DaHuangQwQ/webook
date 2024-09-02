package dao

import (
	"gorm.io/gorm"
	"webook/internal/repository/dao/system"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Article{},
		&PublishedArticle{},
		&Interactive{},
		&UserLikeBiz{},
		&UserCollectionBiz{},
		&system.SysRole{},
		&system.SysAuthRule{},
		&system.SysDept{},
	)
}
