package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/service/system"
	"webook/internal/web"
	"webook/pkg/logger"
)

var _ web.Handler = (*RoleHandler)(nil)

type RoleHandler struct {
	svc system.RoleService
	l   logger.LoggerV1
}

func NewRoleHandler(svc system.RoleService, l logger.LoggerV1) *RoleHandler {
	return &RoleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *RoleHandler) RegisterRoutes(router *gin.Engine) {
	g := router.Group("/role")
	g.GET("/list", h.GetList)
	g.GET("/getParams", h.GetParams)
	g.POST("/add", h.AddRole)
	g.GET("/get", h.GetRole)
	g.PUT("/edit", h.EditRole)
	g.DELETE("/delete", h.DeleteRole)
}

func (h *RoleHandler) GetList(ctx *gin.Context) {

	RoleName := ctx.Query("roleName")
	Status := ctx.Query("roleStatus")
	pageNum := ctx.Query("pageNum")
	pageSize := ctx.Query("pageSize")
	PageNum, _ := strconv.Atoi(pageNum)
	PageSize, _ := strconv.Atoi(pageSize)
	ui8, _ := strconv.ParseUint(Status, 10, 64)
	search, err := h.svc.GetRoleListSearch(ctx, domain.Role{
		Name:   RoleName,
		Status: uint8(ui8),
	}, api.PageReq{
		PageNum:  PageNum,
		PageSize: PageSize,
	})
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: map[string]any{
			"list": search,
		},
	})
}

// GetRoles ok
func (h *RoleHandler) GetParams(ctx *gin.Context) {
	list, err := h.svc.GetParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("获取角色列表错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: map[string]any{
			"menu": list,
		},
	})
}

// AddRole ok
func (h *RoleHandler) AddRole(ctx *gin.Context) {
	type Req struct {
		Name      string  `json:"name"`
		Status    uint8   `json:"status"`
		ListOrder uint    `json:"listOrder"`
		Remark    string  `json:"remark"`
		MenuIds   []int64 `json:"menuIds"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	err := h.svc.AddRole(ctx, domain.Role{
		Name:      req.Name,
		Status:    req.Status,
		ListOrder: req.ListOrder,
		Remark:    req.Remark,
		MenuIds:   req.MenuIds,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("add role error", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

// GetRole ok
func (h *RoleHandler) GetRole(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	Id, _ := strconv.ParseInt(id, 10, 64)
	roleInfo, err := h.svc.Get(ctx, Id)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	menuIds, err := h.svc.GetFilteredNamedPolicy(ctx, Id)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: map[string]any{
			"role":    roleInfo,
			"menuIds": menuIds,
		},
	})
}

// EditRole ok
func (h *RoleHandler) EditRole(ctx *gin.Context) {
	type Req struct {
		Id        int64   `json:"id" v:"required#角色id必须"`
		Name      string  `json:"name" v:"required#角色名称不能为空"`
		Status    uint8   `json:"status"    `
		ListOrder uint    `json:"listOrder" `
		Remark    string  `json:"remark"    `
		MenuIds   []int64 `json:"menuIds"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	err := h.svc.EditRole(ctx, domain.Role{
		Id:        req.Id,
		Name:      req.Name,
		Status:    req.Status,
		ListOrder: req.ListOrder,
		Remark:    req.Remark,
		MenuIds:   req.MenuIds,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

// DeleteRole ok
func (h *RoleHandler) DeleteRole(ctx *gin.Context) {
	type Req struct {
		Ids []int64 `json:"ids" required`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	err := h.svc.DeleteByIds(ctx, req.Ids)
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  err.Error(),
		})
		h.l.Info("delete role error", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}
