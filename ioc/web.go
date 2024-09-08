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
	"webook/internal/web/system"
	"webook/pkg/ginx/middleware/prometheus"
	"webook/pkg/ginx/middleware/ratelimit"
	"webook/pkg/logger"
	limit "webook/pkg/ratelimit"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler,
	authHdl *system.AuthHandler,
	roleHdl *system.RoleHandler,
	sysUserHdl *system.UserHandler,
	deptHdl *system.DeptHandler,
	monitorHdl *system.MonitorHandler,
	orderHdl *web.OrderHandler,
	recruitmentHdl *web.RecruitmentHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	authHdl.RegisterRoutes(server)
	roleHdl.RegisterRoutes(server)
	sysUserHdl.RegisterRoutes(server)
	deptHdl.RegisterRoutes(server)
	monitorHdl.RegisterRoutes(server)
	orderHdl.RegisterRoutes(server)
	recruitmentHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, l logger.LoggerV1) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "DaHuang",
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	return []gin.HandlerFunc{
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token"},
			//AllowHeaders: []string{"content-type"},
			AllowMethods: []string{"POST", "PUT", "DELETE", "OPTIONS", "GET"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "ceitlab.top")
			},
			MaxAge: 12 * time.Hour,
		}),
		ratelimit.NewBuilder(limit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)).Build(),
		middleware.NewLoginJwtMiddleware(l).Build(jwt.NewRedisJWTHandler(redisClient)),
	}
}
