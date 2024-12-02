package dao

import (
	"github.com/DaHuangQwQ/webook/article/repository/dao"
	dao2 "github.com/DaHuangQwQ/webook/internal/cronjob/repository/dao"
	dao3 "github.com/DaHuangQwQ/webook/user/repository/dao"
	"gorm.io/gorm"
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
