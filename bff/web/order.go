package web

import (
	"github.com/gin-gonic/gin"
	"webook/bff/api"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/pkg/ginx"
)

var _ Handler = (*OrderHandler)(nil)

type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{
		svc: svc,
	}
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.GET(ginx.Warp[api.OrderListReq](h.List))
	router.POST(ginx.Warp[api.OrderAddReq](h.Add))
	router.POST(ginx.Warp[api.OrderEditReq](h.Edit))
	router.POST(ginx.Warp[api.OrderDeleteReq](h.Delete))
}

func (h *OrderHandler) List(ctx *gin.Context, req api.OrderListReq) (ginx.Result, error) {
	list, err := h.svc.List(ctx)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: list,
	}, nil
}

func (h *OrderHandler) Add(ctx *gin.Context, req api.OrderAddReq) (ginx.Result, error) {
	err := h.svc.Add(ctx, dao.Order{
		Title:    req.Title,
		Content:  req.Content,
		Author:   req.Author,
		Salary:   req.Salary,
		Duration: req.Duration,
		Address:  req.Address,
		Phone:    req.Phone,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
	}, nil
}

func (h *OrderHandler) Edit(ctx *gin.Context, req api.OrderEditReq) (ginx.Result, error) {
	err := h.svc.Edit(ctx, dao.Order{
		Title:    req.Title,
		Content:  req.Content,
		Author:   req.Author,
		Salary:   req.Salary,
		Duration: req.Duration,
		Address:  req.Address,
		Phone:    req.Phone,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
	}, nil
}

func (h *OrderHandler) Delete(ctx *gin.Context, req api.OrderDeleteReq) (ginx.Result, error) {
	err := h.svc.DeleteByIds(ctx, req.Ids)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
	}, nil
}
