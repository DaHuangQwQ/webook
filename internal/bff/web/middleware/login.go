package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gob.Register(time.Now())
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")

		if sess.Get("userId") == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		now := time.Now().UnixMilli()

		sess.Options(sessions.Options{
			MaxAge: 60,
		})

		if updateTime == nil {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				return
			}
		}
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				return
			}
			return
		}

	}
}
