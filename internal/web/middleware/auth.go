package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
	"webook/internal/domain"
	"webook/internal/service/system"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

type AuthMiddleware struct {
	svc     system.UserService
	authSvc system.AuthService

	casbin casbin.IEnforcer
	l      logger.LoggerV1
}

func NewAuthMiddleware(svc system.UserService, authSvc system.AuthService, casbin casbin.IEnforcer, l logger.LoggerV1) *AuthMiddleware {
	return &AuthMiddleware{
		svc:     svc,
		authSvc: authSvc,
		casbin:  casbin,
		l:       l,
	}
}

func (m *AuthMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Config struct {
			ExcludePaths []string `yaml:"excludePaths"`
		}
		var config Config
		if err := viper.UnmarshalKey("system", &config); err != nil {
			panic(err)
		}
		//获取登陆用户id
		userClaims := ctx.MustGet("claims")
		accessParams := ctx.QueryArray("accessParams")
		accessParamsStr := ""
		if len(accessParams) > 0 && accessParams[0] != "undefined" {
			accessParamsStr = "?" + strings.Join(accessParams, "&")
		}
		url := strings.TrimLeft(ctx.Request.URL.Path, "/") + accessParamsStr

		if userClaims == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user, ok := userClaims.(ijwt.UserClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if m.svc.NotCheckAuthAdminIds(ctx, user.Uid) {
			ctx.Next()
			return
		}

		//获取地址对应的菜单id
		menuList, err := m.authSvc.GetMenuList(ctx)
		if err != nil {
			m.l.Info("请求数据失败", logger.Error(err))
			ctx.JSON(http.StatusUnauthorized, ginx.Result{
				Code: 4,
				Msg:  "请求数据失败",
			})
			return
		}
		var menu *domain.SysAuthRule
		for _, m := range menuList {
			var ms string
			if strings.Index(m.Name, "?") == -1 {
				ms = m.Name
			} else {
				ms = m.Name[0:strings.Index(m.Name, "?")]
			}
			if m.Name == url || ms == url {
				menu = m
				break
			}
		}
		//只验证存在数据库中的规则
		if menu != nil {
			//若是不登录能访问的接口则不判断权限
			for _, p := range config.ExcludePaths {
				if menu.Name == strings.TrimLeft(p, "/") {
					ctx.Next()
					return
				}
			}
			//若存在不需要验证的条件则跳过
			if menu.Condition == "nocheck" {
				ctx.Next()
				return
			}
			menuId := menu.Id
			//菜单没存数据库不验证权限
			if menuId != 0 {
				//判断权限操作
				if err != nil {
					m.l.Info("获取权限失败", logger.Error(err))
					ctx.JSON(http.StatusUnauthorized, ginx.Result{
						Code: 4,
						Msg:  "获取权限失败",
					})
					return
				}
				// u_ 是 casbin 用户前缀
				hasAccess := false
				hasAccess, err = m.casbin.Enforce(fmt.Sprintf("%s%d", "u_", user.Uid), strconv.Itoa(int(menuId)), "All")
				if err != nil {
					m.l.Info("判断权限失败", logger.Error(err))
					ctx.JSON(http.StatusUnauthorized, ginx.Result{
						Code: 4,
						Msg:  "判断权限失败",
					})
					return
				}
				if !hasAccess {
					ctx.JSON(http.StatusUnauthorized, ginx.Result{
						Code: 4,
						Msg:  "没有访问权限",
					})
					return
				}
			}
		} else if menu == nil && accessParamsStr != "" {
			ctx.JSON(http.StatusUnauthorized, ginx.Result{
				Code: 4,
				Msg:  "没有访问权限",
			})
			return
		}
	}
}
