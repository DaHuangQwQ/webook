package web

import (
	"github.com/gin-gonic/gin"
	"webook/pkg/ginx"
)

type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

type Result = ginx.Result
