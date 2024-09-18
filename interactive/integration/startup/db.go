package startup

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	intrDAO "webook/interactive/repository/dao"
)

var db *gorm.DB

// InitTestDB 测试的话，不用控制并发。等遇到了并发问题再说
func InitTestDB() *gorm.DB {
	if db == nil {
		dsn := "root:root@tcp(localhost:3307)/webook"
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()
			if err == nil {
				break
			}
			log.Println("等待连接 MySQL", err)
		}
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(err)
		}
		err = intrDAO.InitTables(db)
		if err != nil {
			panic(err)
		}
		//db = db.Debug()
	}
	return db
}
