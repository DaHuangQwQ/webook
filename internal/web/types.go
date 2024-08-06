package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(router *gin.Engine)
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
