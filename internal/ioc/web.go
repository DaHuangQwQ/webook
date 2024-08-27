package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/jwt"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
	"webook/pkg/logger"
	limit "webook/pkg/ratelimit"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		ratelimit.NewBuilder(limit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)).Build(),
		middleware.NewLoginJwtMiddleware(l).Build(jwt.NewRedisJWTHandler(redisClient)),
	}
}
