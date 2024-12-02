package main

import (
	"github.com/DaHuangQwQ/webook/pkg/ginx"
	"github.com/DaHuangQwQ/webook/pkg/grpcx"
	"github.com/DaHuangQwQ/webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
	webAdmin  *ginx.Server
}
