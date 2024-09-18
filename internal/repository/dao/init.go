package dao

import (
	"gorm.io/gorm"
	"webook/internal/repository/dao/system"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Order{},
		&Recruitment{},
		&Job{},

		&User{},
		&Article{},
		&PublishedArticle{},
		&system.SysRole{},
		&system.SysAuthRule{},
		&system.SysDept{},
	)
}
