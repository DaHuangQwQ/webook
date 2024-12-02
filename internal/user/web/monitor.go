package web

import (
	"github.com/DaHuangQwQ/webook/pkg/ginx"
	"github.com/DaHuangQwQ/webook/user/service/system"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MonitorHandler struct {
	svc system.MonitorService
}

func NewMonitorHandler(svc system.MonitorService) *MonitorHandler {
	return &MonitorHandler{
		svc: svc,
	}
}

func (h *MonitorHandler) RegisterRoutes(router *gin.Engine) {
	g := router.Group("/monitor")
	g.GET("/server", h.Server)
}

func (h *MonitorHandler) Server(ctx *gin.Context) {
	list, err := h.svc.List(ctx)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: list,
	})
}
