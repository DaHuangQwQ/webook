package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

type LoginJwtMiddleware struct {
	paths []string
	l     logger.LoggerV1
}

func NewLoginJwtMiddleware(l logger.LoggerV1) *LoginJwtMiddleware {
	return &LoginJwtMiddleware{
		paths: []string{
			"/users/signup",
			"/users/login",
			"/users/login_sms",
			"/users/login_sms/code/send",
			"/hello",
			"/oauth2/wechat/authurl",
			"/oauth2/wechat/callback",
			"/oauth2/wechat/code2session",
		},
		l: l,
	}
}

func (m *LoginJwtMiddleware) Build(jwtHandler ijwt.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range m.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		tokenString := jwtHandler.ExtractToken(c)
		if tokenString == "" {
			m.l.Info("token为空")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := &ijwt.UserClaims{}
		parseToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			m.l.Info("token解析失败")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if parseToken == nil || !parseToken.Valid {
			// 在这里发现 access_token 过期了，生成一个新的 access_token
			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			m.l.Info("token过期")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != c.Request.UserAgent() {
			// token 被窃取
			// 监控
			m.l.Info("userAgent不同")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = jwtHandler.CheckSession(c, claims.Ssid)
		if err != nil {
			m.l.Info("检查token失败")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("claims", claims)
	}
}
