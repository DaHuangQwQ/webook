package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/service/system"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
)

var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc system.UserService
	l   logger.LoggerV1
}

func NewUserHandler(svc system.UserService, l logger.LoggerV1) *UserHandler {
	return &UserHandler{
		svc: svc,
		l:   l,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	g := router.Group("/user")
	g.GET("/getUserMenus", h.GetUserMenus)
	g.GET("/list", h.List)
	g.GET("/params")
	g.POST("/add")
	g.PUT("/edit")
	g.GET("/getEdit")
	g.PUT("/resetPwd")
	g.PUT("/setStatus")
	g.DELETE("/delete")
	g.GET("/getUsers")
}

func (h *UserHandler) GetUserMenus(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(ijwt.UserClaims)
	userMenus, permissions, err := h.svc.GetAdminRules(ctx, claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("GetAdminRules 错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
		Data: map[string]any{
			"menuList":    userMenus,
			"permissions": permissions,
		},
	})
}

func (h *UserHandler) List(ctx *gin.Context) {
	//pageNum := ctx.Query("pageNum")
	//pageSize := ctx.Query("pageSize")
	//deptId := ctx.Query("deptId")
	//mobile := ctx.Query("mobile")
	//status := ctx.Query("status")
	//keyWords := ctx.Query("keyWords")

}
