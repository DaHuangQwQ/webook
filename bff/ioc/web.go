package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
	"webook/bff/web"
	ijwt "webook/bff/web/jwt"
	"webook/bff/web/middleware"
	"webook/pkg/ginx"
	prometheusx "webook/pkg/ginx/middleware/prometheus"
	"webook/pkg/ginx/middleware/ratelimit"
	"webook/pkg/logger"
	limit "webook/pkg/ratelimit"
)

func InitGinServer(
	redisClient redis.Cmdable,
	l logger.LoggerV1,
	jwtHdl ijwt.Handler,
	user *web.UserHandler,
	article *web.ArticleHandler,
	reward *web.RewardHandler) *ginx.Server {
	pb := &prometheusx.Builder{
		Namespace: "DaHuang",
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	engine := gin.Default()
	engine.Use(
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		corsHdl(),
		timeout(),
		otelgin.Middleware("webook"),
		ratelimit.NewBuilder(limit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)).Build(),
		middleware.NewJWTLoginMiddlewareBuilder(jwtHdl).Build())

	user.RegisterRoutes(engine)
	article.RegisterRoutes(engine)
	reward.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "daming_geektime",
		Subsystem: "webook_bff",
		Name:      "http",
	})
	ginx.NewWarpLogger(l)
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}

func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Request.Context().Deadline()
		if !ok {
			// 强制给一个超时，省得我前端调试等得不耐烦
			newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
			defer cancel()
			ctx.Request = ctx.Request.Clone(newCtx)
		}
		ctx.Next()
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 你不加这个，前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
