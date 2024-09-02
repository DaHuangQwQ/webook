package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webook/internal/api"
	"webook/internal/service/system"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
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
	g.GET("/list", ginx.Warp[api.UserSearchReq](h.List))
	g.GET("/params", ginx.Warp[api.UserGetParamsReq](h.GetParams))
	g.POST("/add", ginx.Warp[api.SetUserReq](h.Add))
	g.PUT("/edit")
	g.GET("/getEdit", ginx.Warp[api.UserGetEditReq](h.GetEdit))
	g.PUT("/resetPwd")
	g.PUT("/setStatus")
	g.DELETE("/delete", ginx.Warp[api.UserDeleteReq](h.Delete))
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

func (h *UserHandler) List(ctx *gin.Context, req api.UserSearchReq) (ginx.Result, error) {
	i64, _ := strconv.ParseInt(ctx.Query("pageNum"), 10, 64)
	req.PageNum = int(i64)
	i64, _ = strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	req.PageSize = int(i64)
	req.DeptId = ctx.Query("deptId")
	req.Mobile = ctx.Query("mobile")
	req.Status = ctx.Query("status")
	req.KeyWords = ctx.Query("keyWords")
	search, err := h.svc.GetUserSearch(ctx, req)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Msg:  "ok",
		Data: search,
	}, nil
}

func (h *UserHandler) Add(ctx *gin.Context, req api.SetUserReq) (ginx.Result, error) {
	err := h.svc.Add(ctx, req)
	if err != nil {
		return ginx.Result{
			Msg: "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "ok",
	}, nil
}

func (h *UserHandler) Delete(ctx *gin.Context, req api.UserDeleteReq) (ginx.Result, error) {
	err := h.svc.Delete(ctx, req.Ids)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  err.Error(),
		}, nil
	}
	return ginx.Result{
		Msg: "ok",
	}, nil
}

func (h *UserHandler) GetParams(ctx *gin.Context, req api.UserGetParamsReq) (ginx.Result, error) {
	params, err := h.svc.GetParams(ctx)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: params,
	}, nil
}

func (h *UserHandler) GetEdit(ctx *gin.Context, req api.UserGetEditReq) (ginx.Result, error) {
	h.svc.GetEdit(ctx, req.Id)
	return ginx.Result{}, nil
}
