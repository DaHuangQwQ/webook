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

var _ web.Handler = (*DeptHandler)(nil)

type DeptHandler struct {
	svc system.DeptService
	l   logger.LoggerV1
}

func NewDeptHandler(svc system.DeptService, l logger.LoggerV1) *DeptHandler {
	return &DeptHandler{
		svc: svc,
		l:   l,
	}
}

func (h *DeptHandler) RegisterRoutes(router *gin.Engine) {
	g := router.Group("/dept")
	g.GET("/list", h.List)
	g.POST("/add", h.Add)
	g.PUT("/edit", h.Edit)
	g.DELETE("/delete", h.Delete)
	g.GET("/treeSelect", h.TreeSelect)
}

func (h *DeptHandler) List(ctx *gin.Context) {
	deptName := ctx.Query("deptName")
	status := ctx.Query("status")
	pageNum := ctx.Query("pageNum")
	pageSize := ctx.Query("pageSize")
	ui64, _ := strconv.ParseUint(status, 10, 64)
	PageNum, _ := strconv.ParseUint(pageNum, 10, 64)
	PageSize, _ := strconv.ParseUint(pageSize, 10, 64)
	list, err := h.svc.List(ctx.Request.Context(), deptName, uint(ui64), int(PageNum), int(PageSize))
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("List 错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Data: map[string]any{
			"deptList": list,
		},
	})
}

func (h *DeptHandler) Add(ctx *gin.Context) {
	var req domain.SysDept
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *DeptHandler) Edit(ctx *gin.Context) {
	var req domain.SysDept
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}
func (h *DeptHandler) Delete(ctx *gin.Context) {
	type Req struct {
		Id uint64 `json:"id" v:"required#id不能为空"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "ok",
	})
}
func (h *DeptHandler) TreeSelect(ctx *gin.Context) {
	res, err := h.svc.TreeSelect(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Msg:  "ok",
		Data: res,
	})
}
