package web

import (
	"github.com/DaHuangQwQ/gpkg/ginx"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

type Result = ginx.Result
