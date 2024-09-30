package dao

import (
	"gorm.io/gorm"
	"webook/article/repository/dao"
	dao2 "webook/cronjob/repository/dao"
	dao3 "webook/user/repository/dao"
	system2 "webook/user/repository/dao/system"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Order{},
		&Recruitment{},
		&dao2.Job{},

		&dao3.User{},
		&dao.Article{},
		&dao.PublishedArticle{},
		&system2.SysRole{},
		&system2.SysAuthRule{},
		&system2.SysDept{},
	)
}
