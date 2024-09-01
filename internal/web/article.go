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
	svc service.ArticleService

	l logger.LoggerV1
}

func NewArticleHandler(articleSvc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: articleSvc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(router *gin.Engine) {
	server := router.Group("/articles")
	server.POST("/edit", h.Edit)
	server.POST("/publish", h.Publish)
	server.POST("/withdraw", h.Withdraw)
	server.GET("/getlist", h.GetList)
	server.POST("/img_update", h.Img_Update)
	server.GET("/list", ginx.Warp[api.PageReq](h.List))
	server.GET("/detail/:id")
}

func (h *ArticleHandler) List(ctx *gin.Context, req api.PageReq) (Result, error) {
	//println(req.PageSize, req.PageNum)
	parseInt, _ := strconv.ParseInt(ctx.Query("pageNum"), 10, 64)
	req.PageNum = int(parseInt)
	parseInt, _ = strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	req.PageSize = int(parseInt)
	list, err := h.svc.List(ctx, req)
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

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	// new or edit
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "参数错误",
		})
		return
	}
	claims := ctx.MustGet("claims").(ijwt.UserClaims)
	id, err := h.svc.Save(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("文章保存失败", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: id,
	})
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

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	// new or edit
	var req ArticleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "参数错误" + err.Error(),
		})
		return
	}
	claims := ctx.MustGet("claims").(ijwt.UserClaims)
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误" + err.Error(),
		})
		h.l.Info("发布文章失败", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	claims := ctx.MustGet("claims").(ijwt.UserClaims)
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *ArticleHandler) GetList(ctx *gin.Context) {
	articles, err := h.svc.GetList(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("获取列表失败", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: articles,
	})
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
