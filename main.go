package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {
	db, err := gorm.Open(mysql.Open(config.Config.DSN))
	if err != nil {
		panic("数据库链接失败")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Host,
		Password: "",
		DB:       0,
	})

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Origin", "Content-Type", "authorization"},
		//ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	server.Use(middleware.NewLoginJwtMiddleware().Build())

	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(redisClient)
	repo := repository.NewUserRepository(userDao, userCache)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)

	err = server.Run(":8080")
	if err != nil {
		return
	}
}
