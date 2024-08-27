package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) RegisterRoutes(router *gin.Engine) {
	server := router.Group("/articles")
	server.POST("/edit", h.Edit)
}

func (h *ArticleHandler) Edit(c *gin.Context) {
	// new or edit
	c.String(http.StatusOK, "你好")
}
