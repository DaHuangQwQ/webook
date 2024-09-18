package main

import (
	"webook/pkg/grpcx"
	"webook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
}
