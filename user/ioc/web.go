package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
	"webook/pkg/ginx"
	"webook/pkg/ginx/middleware/prometheus"
	"webook/pkg/ginx/middleware/ratelimit"
	limit "webook/pkg/ratelimit"
	"webook/user/web"
)

func InitWebServer(mdls []gin.HandlerFunc,
	authHdl *web.AuthHandler,
	roleHdl *web.RoleHandler,
	sysUserHdl *web.UserHandler,
	deptHdl *web.DeptHandler,
	monitorHdl *web.MonitorHandler,
) *ginx.Server {
	server := gin.Default()
	server.Use(mdls...)
	authHdl.RegisterRoutes(server)
	roleHdl.RegisterRoutes(server)
	sysUserHdl.RegisterRoutes(server)
	deptHdl.RegisterRoutes(server)
	monitorHdl.RegisterRoutes(server)
	return &ginx.Server{
		Engine: server,
		Addr:   ":8079",
	}
}

func InitGinMiddlewares(
	redisClient redis.Cmdable,
) []gin.HandlerFunc {
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
		otelgin.Middleware("webook"),
		ratelimit.NewBuilder(limit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)).Build(),
	}
}
