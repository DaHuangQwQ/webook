package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webook/bff/api"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
	"webook/user/service/system"
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
	router.GET(ginx.Warp[api.UserSearchReq](h.List))
	router.GET(ginx.Warp[api.UserGetParamsReq](h.GetParams))
	router.POST(ginx.Warp[api.UserAddReq](h.Add))
	router.PUT(ginx.Warp[api.UserEditReq](h.Edit))
	router.GET(ginx.Warp[api.UserGetEditReq](h.GetEdit))
	g.PUT("/resetPwd")
	router.POST(ginx.Warp[api.UserStatusReq](h.SetStatus))
	router.DELETE(ginx.Warp[api.UserDeleteReq](h.Delete))
	g.GET("/getUsers")
}

func (h *UserHandler) GetUserMenus(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(ijwt.UserClaims)
	userMenus, permissions, err := h.svc.GetAdminRules(ctx.Request.Context(), claims.Uid)
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
	search, err := h.svc.GetUserSearch(ctx.Request.Context(), req)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{
		Msg:  "ok",
		Data: search,
	}, nil
}

func (h *UserHandler) Add(ctx *gin.Context, req api.UserAddReq) (ginx.Result, error) {
	err := h.svc.Add(ctx.Request.Context(), req)
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
	err := h.svc.Delete(ctx.Request.Context(), req.Ids)
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
	params, err := h.svc.GetParams(ctx.Request.Context())
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
	userId := ctx.Query("id")
	parsedUserId, err := strconv.ParseInt(userId, 10, 64)
	req.Id = uint64(parsedUserId)
	UserAndRoleIds, err := h.svc.GetEdit(ctx.Request.Context(), req.Id)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: UserAndRoleIds,
	}, nil
}

func (h *UserHandler) Edit(ctx *gin.Context, req api.UserEditReq) (ginx.Result, error) {
	err := h.svc.Edit(ctx.Request.Context(), req)
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

func (h *UserHandler) SetStatus(ctx *gin.Context, req api.UserStatusReq) (ginx.Result, error) {
	err := h.svc.ChangeUserStatus(ctx.Request.Context(), req.Id, req.UserStatus)
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
