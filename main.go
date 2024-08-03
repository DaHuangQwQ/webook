package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
)
import "gorm.io/driver/mysql"

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/test"))
	if err != nil {
		return
	}

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	//store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16, "tcp", "localhost", "", []byte("123"), []byte("123"))
	//if err != nil {
	//	panic(err)
	//}

	//login := &middleware.LoginMiddlewareBuilder{}
	//server.Use(sessions.Sessions("mysession", store), login.CheckLogin())

	server.Use(middleware.NewLoginJwtMiddleware().Build())

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	userDao := dao.NewUserDao(db)
	repo := repository.NewUserRepository(userDao)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)

	err = server.Run(":8080")
	if err != nil {
		return
	}
}
