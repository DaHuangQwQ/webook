package main

import (
	"webook/pkg/ginx"
	"webook/pkg/grpcx"
	"webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
	webAdmin  *ginx.Server
}
