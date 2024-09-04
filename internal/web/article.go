package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

var _ Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc      service.ArticleService
	interSvc service.InteractiveService

	biz string
	l   logger.LoggerV1
}

func NewArticleHandler(articleSvc service.ArticleService, l logger.LoggerV1, interSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:      articleSvc,
		l:        l,
		interSvc: interSvc,
		biz:      "articles",
	}
}

func (h *ArticleHandler) RegisterRoutes(router *gin.Engine) {
	server := router.Group("/articles")
	router.POST(ginx.WarpWithToken[api.ArticleEditReq](h.Edit))
	router.POST(ginx.WarpWithToken[api.ArticlePublishReq](h.Publish))
	router.POST(ginx.WarpWithToken[api.ArticleWithdrawReq](h.Withdraw))
	router.GET(ginx.Warp[api.ArticleGetListReq](h.GetList))
	server.POST("/img_update", h.Img_Update)
	router.GET(ginx.Warp[api.GetListReq](h.List))
	server.GET("/detail/:id")
	router.POST(ginx.WarpWithToken[api.GetPubArticleDetail](h.GetPubArticle))
	router.POST(ginx.WarpWithToken[api.LikeReq](h.Like))
}

func (h *ArticleHandler) List(ctx *gin.Context, req api.GetListReq) (Result, error) {
	//println(req.PageSize, req.PageNum)
	parseInt, _ := strconv.ParseInt(ctx.Query("pageNum"), 10, 64)
	req.PageNum = int(parseInt)
	parseInt, _ = strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	req.PageSize = int(parseInt)
	list, err := h.svc.List(ctx, api.PageReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		Msg:  "ok",
		Data: list,
	}, nil
}

func (h *ArticleHandler) Edit(ctx *gin.Context, req api.ArticleEditReq, u ijwt.UserClaims) (ginx.Result, error) {
	// new or edit
	id, err := h.svc.Save(ctx, req.ToDomain(u.Uid))
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: id,
	}, nil
}

func (h *ArticleHandler) Img_Update(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	files := form.File["file"]
	var ossAdrs []string
	for _, file := range files {
		if file == nil || err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "图片上传失败",
			})
			return
		}
		if file.Size > 10*1024*1024 {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "文件不应该超过10MB",
			})
			return
		}
		fileType := strings.Split(file.Header.Get("Content-Type"), "/")[1]
		fileOpen, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		ossAdress, err := h.svc.Img_Update(ctx, fileOpen, fileType)
		ossAdrs = append(ossAdrs, ossAdress)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("AvatarUpdate错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: ossAdrs,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context, req api.ArticlePublishReq, u ijwt.UserClaims) (ginx.Result, error) {
	// new or edit
	id, err := h.svc.Publish(ctx, req.ToDomain(u.Uid))
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: id,
	}, nil
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context, req api.ArticleWithdrawReq, u ijwt.UserClaims) (ginx.Result, error) {
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: u.Uid,
		},
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
	}, nil
}

func (h *ArticleHandler) GetList(ctx *gin.Context, req api.ArticleGetListReq) (Result, error) {
	articles, err := h.svc.GetList(ctx)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: articles,
	}, nil
}

func (h *ArticleHandler) Like(ctx *gin.Context, req api.LikeReq, u ijwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		err = h.interSvc.Like(ctx, h.biz, req.ArticleID, u.Uid)
	} else {
		err = h.interSvc.CancelLike(ctx, h.biz, req.ArticleID, u.Uid)
	}
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
	}, nil
}

func (h *ArticleHandler) GetPubArticle(ctx *gin.Context, req api.GetPubArticleDetail, u ijwt.UserClaims) (ginx.Result, error) {
	article, err := h.svc.GetPublishedById(ctx, u.Uid, req.ArticleId)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: article,
	}, nil
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	ImgUrl  string `json:"imgUrl"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (art *ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: uid,
		},
		ImgUrl: art.ImgUrl,
		Type:   art.Type,
	}
}
