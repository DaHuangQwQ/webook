package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
