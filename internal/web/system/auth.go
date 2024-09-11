package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webook/internal/domain"
	"webook/internal/service/system"
	"webook/internal/web"
	"webook/pkg/logger"
)

var _ web.Handler = (*AuthHandler)(nil)

type AuthHandler struct {
	svc     system.AuthService
	roleSvc system.RoleService
	l       logger.LoggerV1
}

func NewAuthHandler(svc system.AuthService, l logger.LoggerV1, roleSvc system.RoleService) *AuthHandler {
	return &AuthHandler{
		svc:     svc,
		l:       l,
		roleSvc: roleSvc,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	g := router.Group("/menu")
	g.GET("/list", h.List)
	g.POST("/add", h.Add)
	g.GET("/getParams", h.GetParams)
	g.GET("/get", h.Get)
	g.PUT("/update", h.Update)
	g.DELETE("/delete", h.Delete)
}

func (h *AuthHandler) List(ctx *gin.Context) {
	type Req struct {
		Title     string `json:"menuName" `
		Component string `json:"component"`
	}
	authList, err := h.svc.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("getList 错误", logger.Field{"err", err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
		Data: map[string]any{
			"rules": authList,
		},
	})
}

func (h *AuthHandler) Add(ctx *gin.Context) {
	var req domain.SysAuthRule
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误" + err.Error(),
		})
		return
	}
	req.Id = 0
	err := h.svc.Add(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("Auth添加失败", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *AuthHandler) Update(ctx *gin.Context) {
	var req domain.SysAuthRule
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	err := h.svc.Update(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("Auth添加失败", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *AuthHandler) Delete(ctx *gin.Context) {
	type Req struct {
		Ids []int64 `json:"ids" `
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	err := h.svc.Delete(ctx.Request.Context(), req.Ids)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("Delete 错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *AuthHandler) GetParams(ctx *gin.Context) {
	menuList, err := h.svc.GetIsMenuList(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("GetIsMenuList 错误", logger.Field{Key: "err", Val: err})
		return
	}
	roleList, err := h.roleSvc.GetRoleList(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("GetIsMenuList 错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
		Data: map[string]any{
			"menus": menuList,
			"roles": roleList,
		},
	})
}

func (h *AuthHandler) Get(ctx *gin.Context) {
	id := ctx.Query("id")
	i64, _ := strconv.ParseInt(id, 10, 64)
	auth, err := h.svc.Get(ctx.Request.Context(), i64)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("Get 错误", logger.Field{Key: "err", Val: err})
		return
	}
	roles, err := h.svc.GetMenuRoles(ctx.Request.Context(), i64)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("GetMenuRoles 错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
		Data: map[string]any{
			"rule":    auth,
			"roleIds": roles,
		},
	})
}
