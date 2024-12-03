package main

import (
	"github.com/DaHuangQwQ/gpkg/ginx"
	"github.com/DaHuangQwQ/gpkg/grpcx"
)

type App struct {
	server *grpcx.Server
	web    *ginx.Server
}
