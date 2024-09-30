package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/pkg/ginx"
	"webook/user/service/system"
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
