package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"webook/pkg/ginx"
)

//go:generate mockgen -source=./types.go -package=jwtmocks -destination=./mocks/handler.mock.go Handler
type Handler interface {
	ClearToken(ctx *gin.Context) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, ssid string, uid int64) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractTokenString(ctx *gin.Context) string
}

type RefreshClaims struct {
	Id   int64
	Ssid string
	jwt.RegisteredClaims
}

// UserClaims 别名机制，偷个懒，这样就不用修改其它的代码了
type UserClaims = ginx.UserClaims
