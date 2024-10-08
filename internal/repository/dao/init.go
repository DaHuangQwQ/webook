package dao

import (
	"gorm.io/gorm"
	"webook/article/repository/dao"
	dao2 "webook/cronjob/repository/dao"
	dao3 "webook/user/repository/dao"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Order{},
		&Recruitment{},
		&dao2.Job{},

		&dao3.User{},
		&dao.Article{},
		&dao.PublishedArticle{},
		&dao3.SysRole{},
		&dao3.SysAuthRule{},
		&dao3.SysDept{},
	)
}
