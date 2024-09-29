package dao

import (
	"gorm.io/gorm"
	"webook/article/repository/dao"
	"webook/internal/repository/dao/system"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Order{},
		&Recruitment{},
		&Job{},

		&User{},
		&dao.Article{},
		&dao.PublishedArticle{},
		&system.SysRole{},
		&system.SysAuthRule{},
		&system.SysDept{},
	)
}
