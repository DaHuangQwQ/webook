package main

import (
	"github.com/DaHuangQwQ/webook/pkg/ginx"
	"github.com/DaHuangQwQ/webook/pkg/grpcx"
)

type App struct {
	server *grpcx.Server
	web    *ginx.Server
}
