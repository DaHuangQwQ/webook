package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginJwtMiddleware struct {
	paths []string
}

func NewLoginJwtMiddleware() *LoginJwtMiddleware {
	return &LoginJwtMiddleware{
		paths: []string{"/users/signup", "/users/login"},
	}
}

func (m *LoginJwtMiddleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range m.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &web.UserClaims{}
		parseToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("123412341234"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != c.Request.UserAgent() {
			// token 被窃取
			// 监控
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now).Minutes() > 60 {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Second * 30))
			//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, _ := parseToken.SignedString([]byte("123412341234"))
			c.Header("x-jwt-token", tokenStr)
		}
		println(parseToken)
		c.Set("claims", claims)
	}
}
