package main

import (
	"webook/pkg/ginx"
	"webook/pkg/grpcx"
)

type App struct {
	server   *grpcx.Server
	webAdmin *ginx.Server
}
